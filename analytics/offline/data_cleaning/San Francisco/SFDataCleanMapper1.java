import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.NullWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Mapper;
import java.io.IOException;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.text.ParseException;


public class SFDataCleanMapper1 extends Mapper<LongWritable, Text, NullWritable, Text>  {

    @Override
    public void map(LongWritable key, Text value, Context context) throws IOException, InterruptedException {

        if (key.get() == 0) {
            return;
        }

        // Split the input line by tab
        String[] fields = value.toString().split("\t");

        // Build a string with selected columns separated by tabs
        StringBuilder selectedFields = new StringBuilder();
        //longitude
        if(!fields[12].equals("-120.5")){
                selectedFields.append(fields[11] + "\t");
    	}
	else return;
        //latitude
        if(!fields[13].equals("90")){
	        selectedFields.append(fields[12] + "\t");
	}
	else return;

        SimpleDateFormat dateFormat = new SimpleDateFormat("MM/dd/yyyy hh:mm");
        try {
            Date date = dateFormat.parse(fields[6] + " " + fields[7]);
            long unixTimestamp = date.getTime() / 1000;
            selectedFields.append(unixTimestamp);
        } catch (ParseException e) {
            selectedFields.append(e.getMessage());
        }
        selectedFields.append("\t");

        selectedFields.append(fields[4]);

        // Emit the selected fields with null as key
        context.write(NullWritable.get(), new Text(selectedFields.toString().trim()));
    }
}
