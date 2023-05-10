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
import java.time.ZoneOffset;
import java.text.ParseException;


public class LAMapper extends Mapper<LongWritable, Text, NullWritable, Text>  {

    @Override
    public void map(LongWritable key, Text value, Context context) throws IOException, InterruptedException {

        if (key.get() == 0) {
            return;
        }

        // Split the TSV file
        String[] line = value.toString().split("\t");

        // Use StringBuilder object to append extracted columns
        StringBuilder extractedColumns = new StringBuilder();

        String longitude = line[line.length - 1];
        // Trim trailing or leading whitespaces
        longitude = longitude.trim();
        // Data Dictionary specifies that lat and long can be empty or 0 for missing data
        if (!longitude.isEmpty() && !longitude.equals("0")) {
            extractedColumns.append(longitude);
            extractedColumns.append("\t");
        } else {
            // Do not take this row
            return;
        }

        String latitude = line[line.length - 2];
        latitude = latitude.trim();
        if (!latitude.isEmpty() && !latitude.equals("0")) {
            extractedColumns.append(latitude);
            extractedColumns.append("\t");
        } else {
            // Do not take this row
            return;
        }

        // Define the format types
        SimpleDateFormat dateFormat = new SimpleDateFormat("MM/dd/yyyy hh:mm:ss a");
        DateTimeFormatter timeFormatter = DateTimeFormatter.ofPattern("HHmm");

        try {
            Date date = dateFormat.parse(line[2]);
            String dateString = new SimpleDateFormat("yyyy-MM-dd").format(date);
            LocalDate datePart = LocalDate.parse(dateString);
            LocalTime localTime = LocalTime.parse(line[3], timeFormatter);
            LocalDateTime dt = LocalDateTime.of(datePart, localTime);

            DateTimeFormatter fullDateFormatter = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");
            String formattedDateTime = dt.format(fullDateFormatter);
            LocalDateTime dateTime = LocalDateTime.parse(formattedDateTime, fullDateFormatter);

            // Calculate final unix timestamp
            long unixTimestamp = dateTime.atZone(ZoneOffset.UTC).toInstant().getEpochSecond();
            extractedColumns.append(Long.toString(unixTimestamp));
            extractedColumns.append("\t");
        } catch (ParseException e) {
            // Initially I was returning here
            // But then I modified the code to catch the exception to figure out my error
            extractedColumns.append(e.getMessage());
        }
        extractedColumns.append(line[9]);
        context.write(NullWritable.get(), new Text(extractedColumns.toString().trim()));
    }
}
