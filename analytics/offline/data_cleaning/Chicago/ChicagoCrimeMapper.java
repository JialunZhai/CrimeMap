import java.io.IOException;

import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.io.NullWritable;
import org.apache.hadoop.mapreduce.Mapper;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.text.ParseException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ChicagoCrimeMapper extends Mapper<LongWritable, Text, Text, NullWritable>{
    static int date = 2;
    static int type = 5;
    static int arrest = 8;

    @Override
    public void map(LongWritable key, Text value, Context context) throws IOException, InterruptedException {
        String line = value.toString();
        String[] words = line.split("\t");

        int len = words.length;
        // if(len >0 && (words[0].substring(1, 3).equals("ID"))) return;
        if(len == 0) return;
        String output = ""; 

        Pattern pattern = Pattern.compile("\\((-?\\d+\\.\\d+)?,\\s*(-?\\d+\\.\\d+)\\)");
        String location = words[len-1];
        Matcher matcher = pattern.matcher(location);

        if (matcher.matches()) {
            String latitude = matcher.group(1);
            Float.parseFloat(latitude);
            String longitude = matcher.group(2);
            Float.parseFloat(longitude);
            output += latitude + "\t" + longitude;
        } else {
            return;
        }

        SimpleDateFormat dateFormat = new SimpleDateFormat("M/d/yyyy h:mm:ss a");
        try {
            Date time = dateFormat.parse(words[date]);
            long unixTimestamp = time.getTime() / 1000L;
            output += "\t" + Long.toString(unixTimestamp);
        } catch (ParseException e) {
            return;
        }

        output += "\t" + words[type] + "," + words[arrest];

        context.write(new Text(output), NullWritable.get());
    }
}
