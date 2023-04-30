package hbase_client

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	env_interface "github.com/jialunzhai/crimemap/analytics/online/server/enviroment"
	"github.com/jialunzhai/crimemap/analytics/online/server/interfaces"
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
	longitudeQualifier   = "x"
	laitudeQualifier     = "y"
	timeQualifier        = "t"
	descriptionQualifier = "d"
	yyyyMMddTHHmmss      = "2006-01-02T15:04:05"
)

func Register(env env_interface.Env) error {
	config := env.GetConfig()
	if config == nil || config.Database.Address == "" || config.Database.Namespace == "" || config.Database.Table == "" {
		return errors.New("HBase client not configured")
	}
	c, err := NewHBaseClient(env, config.Database.Address, config.Database.Namespace, config.Database.Table)
	if err != nil {
		return err
	}
	env.SetDatabaseClient(c)
	return nil
}

func NewHBaseClient(env env_interface.Env, zkquorum, namespace, table string) (*HBaseClient, error) {
	client := gohbase.NewClient(zkquorum)
	return &HBaseClient{
		env:    env,
		client: client,
		table:  fmt.Sprintf("%s:%s", namespace, table),
	}, nil
}

func (c *HBaseClient) Conn(ctx context.Context) error {
	/*
		// TODO: move the table and rowKey into config
		getRequest, err := hrpc.NewGetStr(ctx, "group4:test", "dr5rugb9rwjj1970-01-20T05:54:52")
		if err != nil {
			return fmt.Errorf("hrpc error: %v\n", err)
		}
		getRsp, err := c.client.Get(getRequest)
		if err != nil {
			return fmt.Errorf("get HBase response failed: %v\n", err)
		}
		_ = getRsp
		//fmt.Printf("DEBUG: HBase response: %v\n", getRsp)
	*/
	return nil
}

func (c *HBaseClient) GetCrimes(ctx context.Context, minX, maxX, minY, maxY float64, minT, maxT int64) ([]*interfaces.Crime, error) {
	crimes := make([]*interfaces.Crime, 64)

	minHash := geohash.Encode(minY, minX, maxPrecision)
	maxHash := geohash.Encode(maxY, maxX, maxPrecision)
	prefixRowKey := longestCommonPrefix(minHash, maxHash)
	prefixRowKeyFilter := filter.NewPrefixFilter([]byte(prefixRowKey))

	minNormalizedX := normalizeCoordinate(minX, -180.0)
	maxNormalizedX := normalizeCoordinate(maxX, -180.0)
	rangeColXFilter := filter.NewColumnRangeFilter([]byte(minNormalizedX), []byte(maxNormalizedX), true, true)

	minNormalizedY := normalizeCoordinate(minY, -90.0)
	maxNormalizedY := normalizeCoordinate(maxY, -90.0)
	rangeColYFilter := filter.NewColumnRangeFilter([]byte(minNormalizedY), []byte(maxNormalizedY), true, true)

	minNormalizedT := normalizeTime(minT)
	maxNormalizedT := normalizeTime(maxT)
	rangeColTFilter := filter.NewColumnRangeFilter([]byte(minNormalizedT), []byte(maxNormalizedT), true, true)

	scanRequest, err := hrpc.NewScanStr(ctx, "group4_rbda_nyu_edu:crimes",
		hrpc.Filters(prefixRowKeyFilter), hrpc.Filters(rangeColXFilter), hrpc.Filters(rangeColYFilter), hrpc.Filters(rangeColTFilter))
	scanRsp := c.client.Scan(scanRequest)

	var result *hrpc.Result
	for {
		result, err = scanRsp.Next()
		if err != nil {
			break
		}
		var crime interfaces.Crime
		for _, cell := range result.Cells {
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
			// TODO: remove these stmts
			// --- begin ---
			if crime.Longitude < minX || crime.Longitude > maxX || crime.Latitude < minY || crime.Latitude > maxY || crime.Time < minT || crime.Time > maxT {
				continue
			}
			// --- end ---
			crimes = append(crimes, &crime)
		}
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
			return s1[:i-1]
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

func normalizeTime(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(yyyyMMddTHHmmss)
}

func denormalizeCoordinate(normalizedStr string, minValue float64) float64 {
	// TODO: Don't ignore this error
	value, _ := strconv.ParseFloat(normalizedStr, 64)
	return value/1e6 + minValue
}

func denormalizeTime(normalizeTime string) int64 {
	// TODO: Don't ignore this error
	date, _ := time.Parse(yyyyMMddTHHmmss, normalizeTime)
	return date.Unix()
}
