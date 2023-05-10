import java.io.IOException;
import java.text.SimpleDateFormat;
import java.util.Date;
import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.NullWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Mapper;

public class DataCleaningNYCCurrentMapper
    extends Mapper<LongWritable, Text, NullWritable, Text> {

  private static final String tableHead = "CMPLNT_NUM	ADDR_PCT_CD	BORO_NM	CMPLNT_FR_DT	CMPLNT_FR_TM	CMPLNT_TO_DT	CMPLNT_TO_TM	CRM_ATPT_CPTD_CD	HADEVELOPT	HOUSING_PSA	JURISDICTION_CODE	JURIS_DESC	KY_CD	LAW_CAT_CD	LOC_OF_OCCUR_DESC	OFNS_DESC	PARKS_NM	PATROL_BORO	PD_CD	PD_DESC	PREM_TYP_DESC	RPT_DT	STATION_NAME	SUSP_AGE_GROUP	SUSP_RACE	SUSP_SEX	TRANSIT_DISTRICT	VIC_AGE_GROUP	VIC_RACE	VIC_SEX	X_COORD_CD	Y_COORD_CD	Latitude	Longitude	Lat_Lon	New Georeferenced Column";

  @Override
  public void map(LongWritable key, Text value, Context context)
      throws IOException, InterruptedException {
    String line = value.toString();
    if (line.length() == 0 || line.equals(tableHead))
      return;

    String[] fields = line.split("\t");
    if (fields.length != 36) {
      context.getCounter("BadRecords", "BAD_FIELDS_NUM").increment(1);
      return;
    }

    // extract raw values of fields
    String day = fields[3].trim();
    String time = fields[4].trim();
    String offenseLevel = fields[13].trim();
    String offenseDesc = fields[15].trim();
    String detailedDesc = fields[19].trim();
    String latitude = fields[32].trim();
    String longitude = fields[33].trim();

    // get the description
    String description;
    if (detailedDesc.length() == 0 || detailedDesc.equals("(null)")) {
      if (offenseDesc.length() == 0 || offenseDesc.equals("(null)")) {
        description = offenseLevel;
      } else {
        description = String.format("%s: %s", offenseLevel, offenseDesc);
      }
    } else {
      description = String.format("%s: %s", offenseLevel, detailedDesc);
    }

    // parse raw values and check whether they are valid
    String timestamp;
    SimpleDateFormat dateFormat = new SimpleDateFormat("MM/dd/yyyy HH:mm:ss");
    try {
      Date date = dateFormat.parse(day + " " + time);
      long unixTimestamp = date.getTime() / 1000;
      timestamp = Long.toString(unixTimestamp);

      // check whether values of latitude and longittude fields are valid float
      float fpLatitude = Float.parseFloat(latitude);
      if (latitude.length() == 0 || Float.parseFloat(latitude) == 0.0
          || longitude.length() == 0 || Float.parseFloat(longitude) == 0.0) {
        context.getCounter("Bad Records", "MISSING_FIELD_VALUE").increment(1);
        return;
      }
    } catch (Exception e) {
      context.getCounter("Bad Records", "BAD_DATA_FORMAT").increment(1);
      return;
    }

    String record = String.format("%s\t%s\t%s\t%s", longitude, latitude, timestamp, description);

    context.write(NullWritable.get(), new Text(record));
  }
}
