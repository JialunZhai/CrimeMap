import ch.hsr.geohash.GeoHash;
import java.io.IOException;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import org.apache.hadoop.hbase.KeyValue;
import org.apache.hadoop.hbase.io.ImmutableBytesWritable;
import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Mapper;

public class DataImportingMapper
    extends Mapper<LongWritable, Text, ImmutableBytesWritable, KeyValue> {
  final static byte[] COL_FAMILY = "e".getBytes();
  ImmutableBytesWritable hKey = new ImmutableBytesWritable();
  KeyValue kv;

  private String normalize(double value, double minValue) {
    double posValue = value + minValue;
    long prefixValue = (long) posValue;
    String prefixStr = String.valueOf(prefixValue);
    // add extra '0' to the front
    String prefixPadding = new String(new char[3 - prefixStr.length()]).replace('\0', '0');
    prefixStr = prefixPadding + prefixStr;

    long suffixValue = (long) ((posValue - (double) prefixValue) * 1e6);
    String suffixStr = String.valueOf(suffixValue);
    // add extra '0' to the back
    String suffixPadding = new String(new char[6 - suffixStr.length()]).replace('\0', '0');
    suffixStr += suffixPadding;

    return prefixStr + suffixStr;
  }

  @Override
  public void map(LongWritable key, Text value, Context context)
      throws IOException, InterruptedException {
    String line = value.toString();
    if (line.length() == 0)
      return;

    String[] fields = line.split("\t");
    if (fields.length != 4) {
      context.getCounter("BadRecords", "BAD_FIELDS_NUM").increment(1);
      return;
    }

    // extract raw values of fields
    String rawLongitude = fields[0];
    String rawLatitude = fields[1];
    String timestamp = fields[2];
    String description = fields[3];

    String normalizedLongitude, normalizedLatitude, time, rowKey;

    // parse raw values and check whether they are valid
    SimpleDateFormat sdf = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss");
    try {
      // check whether values of latitude and longitude fields are valid double
      double longitude = Double.parseDouble(rawLongitude);
      double latitude = Double.parseDouble(rawLongitude);
      // It's impossible that a double value exactly equals to 0.0
      if (longitude == 0.0 || longitude < -180.0 || longitude > 180.0
          || latitude == 0.0 || latitude < -90.0 || longitude > 90.0) {
        context.getCounter("Bad Records", "MISSING_FIELD_VALUE").increment(1);
        return;
      }

      normalizedLongitude = normalize(longitude);
      normalizedLatitude = normalize(latitude);

      Date date = new Date(timestamp);
      time = sdf.format(date);

      GeoHash geohash = GeoHash.withCharacterPrecision(latitude, longitude, 12);
      rowKey = geohash + time;
    } catch (Exception e) {
      context.getCounter("Bad Records", "BAD_DATA_FORMAT").increment(1);
      return;
    }

    hKey.set(rowKey.getBytes());

    kv = new KeyValue(hKey.get(), COL_FAMILY, HColumnEnum.COL_X.getColumnName(), normalizedLongitude.getBytes());
    context.write(hKey, kv);

    kv = new KeyValue(hKey.get(), COL_FAMILY, HColumnEnum.COL_Y.getColumnName(), normalizedLatitude.getBytes());
    context.write(hKey, kv);

    kv = new KeyValue(hKey.get(), COL_FAMILY, HColumnEnum.COL_T.getColumnName(), time.getBytes());
    context.write(hKey, kv);

    if (description.length() != 0) {
      kv = new KeyValue(hKey.get(), COL_FAMILY, HColumnEnum.COL_D.getColumnName(), description.getBytes());
      context.write(hKey, kv);
    }
  }
}
