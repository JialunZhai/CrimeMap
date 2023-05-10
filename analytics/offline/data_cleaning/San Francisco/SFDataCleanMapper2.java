import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.NullWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Mapper;
import java.io.IOException;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.text.ParseException;

public class SFDataCleanMapper2 extends Mapper<LongWritable, Text, NullWritable, Text>  {

    @Override
    public void map(LongWritable key, Text value, Context context) throws IOException, InterruptedException {

        if (key.get() == 0) {
            return;
        }

        // Split the input line by tab
        String[] fields = value.toString().split("\t");

	//context.write(NullWritable.get(), new Text(String.valueOf(fields.length)));
	if(fields.length != 35) return;
        // Build a string with selected columns separated by tabs
        StringBuilder selectedFields = new StringBuilder();
        //longitude
        if(!fields[25].equals("0")){
                selectedFields.append(fields[25]);
                selectedFields.append("\t");
    	}
	    else return;
            //latitude
        if(!fields[24].equals("0")){
	        selectedFields.append(fields[24]);
            selectedFields.append("\t");
	    }

        SimpleDateFormat dateFormat = new SimpleDateFormat("yyyy/MM/dd hh:mm");
        try {
            Date date = dateFormat.parse(fields[1] + " " + fields[2]);
            long unixTimestamp = date.getTime() / 1000;
            selectedFields.append(unixTimestamp);
        } catch (ParseException e) {
            selectedFields.append(e.getMessage());
        }
        selectedFields.append("\t");

        selectedFields.append(fields[11] + "," + fields[14]);

        // Emit the selected fields with null as key
        context.write(NullWritable.get(), new Text(selectedFields.toString().trim()));
    }
}
