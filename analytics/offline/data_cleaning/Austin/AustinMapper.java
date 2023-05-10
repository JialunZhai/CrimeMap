import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.NullWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Mapper;
import java.io.IOException;
import java.text.SimpleDateFormat;
import java.time.LocalTime;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.Date;
import java.time.ZoneId;
import java.text.ParseException;
import java.util.TimeZone;



public class AustinMapper extends Mapper<LongWritable, Text, NullWritable, Text>  {

    @Override
    public void map(LongWritable key, Text value, Context context) throws IOException, InterruptedException {

        if (key.get() == 0) {
            return;
        }

        // Split the TSV file
        String[] line = value.toString().split("\t");

        // Use StringBuilder object to append extracted columns
        StringBuilder extractedColumns = new StringBuilder();

        String longitude = line[line.length - 2];
        // Trim trailing or leading whitespaces
        longitude = longitude.trim();
        // Data Dictionary specifies that lat and long can be empty for missing data
        if (!longitude.isEmpty()) {
            extractedColumns.append(longitude);
            extractedColumns.append("\t");
        } else {
            // Do not take this row
            return;
        }

        String latitude = line[line.length - 3];
        latitude = latitude.trim();
        if (!latitude.isEmpty()) {
            extractedColumns.append(latitude);
            extractedColumns.append("\t");
        } else {
            // Do not take this row
            return;
        }

        // Define the format types
        SimpleDateFormat dateFormat = new SimpleDateFormat("MM/dd/yyyy hh:mm:ss a");

        try {
            Date date = dateFormat.parse(line[4]);
            // Calculate final unix timestamp
            long unixTimestamp = date.getTime() / 1000;
            extractedColumns.append(Long.toString(unixTimestamp));
            extractedColumns.append("\t");
        } catch (ParseException e) {
            extractedColumns.append(e.getMessage());
        }
        extractedColumns.append(line[1]);
        context.write(NullWritable.get(), new Text(extractedColumns.toString().trim()));
    }
}
