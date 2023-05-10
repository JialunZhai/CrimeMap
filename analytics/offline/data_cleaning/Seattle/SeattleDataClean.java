import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.NullWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.lib.input.MultipleInputs;
import org.apache.hadoop.mapreduce.lib.input.TextInputFormat;
import org.apache.hadoop.mapreduce.lib.output.TextOutputFormat;

public class SeattleDataClean {

    public static void main(String[] args) throws Exception {

        if (args.length != 2) {
			System.err.println("Usage: SeattleDataClean <input path> <output path>");
			System.exit(-1);
		}

        // Create a Hadoop job and set the jar
        Configuration conf = new Configuration();
        Job job = Job.getInstance(conf, "seattle data clean");
        job.setJarByClass(SeattleDataClean.class);


        // Set the input and output file formats and paths
        job.setInputFormatClass(TextInputFormat.class);
        job.setOutputFormatClass(TextOutputFormat.class);
        TextOutputFormat.setOutputPath(job, new Path(args[1]));

        // Set the mapper and reducer classes
        MultipleInputs.addInputPath(job, new Path(args[0]), TextInputFormat.class, SeattleDataCleanMapper.class);
        job.setReducerClass(SeattleDataCleanReducer.class);

        // Set the output key and value classes
        job.setOutputKeyClass(NullWritable.class);
        job.setOutputValueClass(Text.class);

        // Submit the job and wait for it to complete
        System.exit(job.waitForCompletion(true) ? 0 : 1);
    }
}
