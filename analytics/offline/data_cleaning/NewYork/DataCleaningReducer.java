import java.io.IOException;
import org.apache.hadoop.io.NullWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Reducer;

public class DataCleaningReducer
    extends Reducer<NullWritable, Text, Text, NullWritable> {

  @Override
  public void reduce(NullWritable key, Iterable<Text> records, Context context)
      throws IOException, InterruptedException {
    float longitude = 0, minLongitude = 180, maxLongitude = -180;
    float latitude = 0, minLatitude = 90, maxLatitude = -90;
    long timestamp = 0, minTimestamp = Long.MAX_VALUE, maxTimestamp = Long.MIN_VALUE;
    for (Text record : records) {
      context.write(record, NullWritable.get());

      // data profiling
      String[] fields = record.toString().split("\t");
      longitude = Float.parseFloat(fields[0]);
      latitude = Float.parseFloat(fields[1]);
      timestamp = Long.parseLong(fields[2]);

      minLongitude = Math.min(minLongitude, longitude);
      maxLongitude = Math.max(maxLongitude, longitude);
      minLatitude = Math.min(minLatitude, latitude);
      maxLatitude = Math.max(maxLatitude, latitude);
      minTimestamp = Math.min(minTimestamp, timestamp);
      maxTimestamp = Math.max(maxTimestamp, timestamp);
    }

    context.getCounter("MinLongitude", "").increment(Math.round(minLongitude));
    context.getCounter("MaxLongitude", "").increment(Math.round(maxLongitude));
    context.getCounter("MinLatitude", "").increment(Math.round(minLatitude));
    context.getCounter("MaxLatitude", "").increment(Math.round(maxLatitude));
    context.getCounter("MinTimestamp", "").increment(minTimestamp);
    context.getCounter("MaxTimestamp", "").increment(maxTimestamp);
  }
}
