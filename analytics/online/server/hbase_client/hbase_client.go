package hbase_client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	env_interface "github.com/jialunzhai/crimemap/analytics/online/server/enviroment"
	"github.com/jialunzhai/crimemap/analytics/online/server/interfaces"
	"github.com/jialunzhai/crimemap/analytics/online/server/metrics"
	"github.com/pierrre/geohash"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/filter"
	"github.com/tsuna/gohbase/hrpc"
)

type HBaseClient struct {
	env    env_interface.Env
	client gohbase.Client
	table  string
}

const (
	maxPrecision         = 12
	eventFamily          = "e"
	longitudeQualifier   = "x"
	laitudeQualifier     = "y"
	timeQualifier        = "t"
	descriptionQualifier = "d"
	yyyyMMddTHHmmss      = "2006-01-02T15:04:05"
)

func Register(env env_interface.Env) error {
	config := env.GetConfig()
	if config == nil || config.Database.Address == "" || config.Database.Table == "" {
		return errors.New("HBase client not configured")
	}
	c, err := NewHBaseClient(env, config.Database.Address, config.Database.Table)
	if err != nil {
		return err
	}
	env.SetDatabaseClient(c)
	return nil
}

func NewHBaseClient(env env_interface.Env, zkquorum, table string) (*HBaseClient, error) {
	log.Printf("table name `%v`\n", table)
	client := gohbase.NewClient(zkquorum)
	return &HBaseClient{
		env:    env,
		client: client,
		table:  table,
	}, nil
}

func (c *HBaseClient) Conn(ctx context.Context) error {
	// c.GetCrimes(ctx, -122.3592, -122.059, 47.5272, 47.5274, 1513799100, 1593799300)
	return nil
}

func (c *HBaseClient) GetCrimes(ctx context.Context, minLongitude, maxLongitude, minLaitude, maxLaitude float64, minTime, maxTime int64) ([]*interfaces.Crime, error) {
	fmt.Printf("DEBUG: GetCrimes received minLong=%v, maxLong=%v, minLa=%v, maxLa=%v, minT=%v, maxT=%v\n",
		minLongitude, maxLongitude, minLaitude, maxLaitude, minTime, maxTime)
	if minLongitude > maxLongitude || minLaitude > maxLaitude || minTime > maxTime {
		return nil, fmt.Errorf("HBase client warnning: bad query arguments for GetCrimes")
	}
	startTime := time.Now()
	crimes := make([]*interfaces.Crime, 0)

	minHash := geohash.Encode(minLaitude, minLongitude, maxPrecision)
	maxHash := geohash.Encode(maxLaitude, maxLongitude, maxPrecision)
	prefixRowKey := longestCommonPrefix(minHash, maxHash)
	prefixRowKeyFilter := filter.NewPrefixFilter([]byte(prefixRowKey))

	minNormalizedX := normalizeCoordinate(minLongitude, -180.0)
	geqMinXFilter := filter.NewSingleColumnValueFilter([]byte(eventFamily), []byte(longitudeQualifier),
		filter.GreaterOrEqual, filter.NewBinaryComparator(filter.NewByteArrayComparable([]byte(minNormalizedX))), true, true)
	maxNormalizedX := normalizeCoordinate(maxLongitude, -180.0)
	leqMaxXFilter := filter.NewSingleColumnValueFilter([]byte(eventFamily), []byte(longitudeQualifier),
		filter.LessOrEqual, filter.NewBinaryComparator(filter.NewByteArrayComparable([]byte(maxNormalizedX))), true, true)

	minNormalizedY := normalizeCoordinate(minLaitude, -90.0)
	geqMinYFilter := filter.NewSingleColumnValueFilter([]byte(eventFamily), []byte(laitudeQualifier),
		filter.GreaterOrEqual, filter.NewBinaryComparator(filter.NewByteArrayComparable([]byte(minNormalizedY))), true, true)
	maxNormalizedY := normalizeCoordinate(maxLaitude, -90.0)
	leqMaxYFilter := filter.NewSingleColumnValueFilter([]byte(eventFamily), []byte(laitudeQualifier),
		filter.LessOrEqual, filter.NewBinaryComparator(filter.NewByteArrayComparable([]byte(maxNormalizedY))), true, true)

	minNormalizedT := normalizeTime(minTime)
	geqMinTFilter := filter.NewSingleColumnValueFilter([]byte(eventFamily), []byte(timeQualifier),
		filter.GreaterOrEqual, filter.NewBinaryComparator(filter.NewByteArrayComparable([]byte(minNormalizedT))), true, true)
	maxNormalizedT := normalizeTime(maxTime)
	leqMaxTFilter := filter.NewSingleColumnValueFilter([]byte(eventFamily), []byte(timeQualifier),
		filter.LessOrEqual, filter.NewBinaryComparator(filter.NewByteArrayComparable([]byte(maxNormalizedT))), true, true)

	filterList := filter.NewList(filter.MustPassAll, prefixRowKeyFilter,
		geqMinTFilter, leqMaxTFilter, geqMinXFilter, leqMaxXFilter, geqMinYFilter, leqMaxYFilter)
	scanRequest, err := hrpc.NewScanStr(ctx, c.table,
		hrpc.Filters(filterList))
	scanRsp := c.client.Scan(scanRequest)

	rowCountHBaseReturned := 0
	rowCountCorrect := 0
	var result *hrpc.Result
	for {
		result, err = scanRsp.Next()
		if err != nil {
			break
		}
		crime := &interfaces.Crime{}
		for _, cell := range result.Cells {
			if string(cell.Family) != eventFamily {
				continue
			}
			switch string(cell.Qualifier) {
			case longitudeQualifier:
				crime.Longitude = denormalizeCoordinate(string(cell.Value), -180.0)
			case laitudeQualifier:
				crime.Latitude = denormalizeCoordinate(string(cell.Value), -90.0)
			case timeQualifier:
				crime.Time = denormalizeTime(string(cell.Value))
			case descriptionQualifier:
				crime.Description = string(cell.Value)
			default:
				// unexpeced qualifier: just skip it to fit changes in schema
			}
		}
		rowCountHBaseReturned++
		// TODO: remove these condition stmts after filter fully tested
		if crime.Longitude < minLongitude || crime.Longitude > maxLongitude || crime.Latitude < minLaitude || crime.Latitude > maxLaitude {
			continue
		}

		if crime.Time < minTime || crime.Time > maxTime {
			continue
		}
		rowCountCorrect++
		//fmt.Printf("%v\n", *crime)
		crimes = append(crimes, crime)
	}
	diff := time.Now().Sub(startTime)
	if rowCountCorrect != rowCountHBaseReturned {
		fmt.Printf("WARNNING: query-prefix length %v out of 12, %v of %v namely %.2f%% the returned records fit the conditions, query costs %v seconds\n",
			len(prefixRowKey), rowCountCorrect, rowCountHBaseReturned, 100*float64(rowCountCorrect)/float64(rowCountHBaseReturned), diff.Seconds())
	} else {
		fmt.Printf("DEBUG: query-prefix length %v out of 12, %v of %v namely %.2f%% the returned records fit the conditions, query costs %v seconds\n",
			len(prefixRowKey), rowCountCorrect, rowCountHBaseReturned, 100*float64(rowCountCorrect)/float64(rowCountHBaseReturned), diff.Seconds())
	}
	return crimes, nil
}

func (c *HBaseClient) Close() error {
	c.client.Close()
	return nil
}

func longestCommonPrefix(s1, s2 string) string {
	for i := 0; i < len(s1) && i < len(s2); i++ {
		if s1[i] != s2[i] {
			return s1[:i]
		}
	}
	if len(s1) < len(s2) {
		return s1
	} else {
		return s2
	}
}

func normalizeCoordinate(value, minValue float64) string {
	posValue := value - minValue
	prefixValue := int64(posValue)
	prefixStr := strconv.FormatInt(prefixValue, 10)
	prefixPadding := strings.Repeat("0", 3-len(prefixStr))
	prefixStr = prefixPadding + prefixStr

	suffixValue := int64((posValue - float64(prefixValue)) * 1e6)
	suffixStr := strconv.FormatInt(suffixValue, 10)
	suffixPadding := strings.Repeat("0", 6-len(suffixStr))
	suffixStr += suffixPadding

	return prefixStr + suffixStr
}

func denormalizeCoordinate(normalizedStr string, minValue float64) float64 {
	value, err := strconv.ParseFloat(normalizedStr, 64)
	if err != nil {
		metrics.InvalidRowsCounter.Inc()
		log.Printf("WARNING: parse `%v` to float64 failed with error: %v\n", normalizedStr, err)
		return 0
	}
	return value/1e6 + minValue
}

func normalizeTime(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(yyyyMMddTHHmmss)
}

func denormalizeTime(normalizeTime string) int64 {
	// TODO: Don't ignore this error
	date, err := time.Parse(yyyyMMddTHHmmss, normalizeTime)
	if err != nil {
		metrics.InvalidRowsCounter.Inc()
		log.Printf("WARNING: parse `%v` to date failed with error: %v\n", normalizeTime, err)
		return 0
	}
	return date.Unix()
}
