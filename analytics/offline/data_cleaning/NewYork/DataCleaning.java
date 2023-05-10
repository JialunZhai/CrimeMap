import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.NullWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;
import org.apache.hadoop.mapreduce.lib.input.MultipleInputs;
import org.apache.hadoop.mapreduce.lib.input.TextInputFormat;

public class DataCleaning {

  public static void main(String[] args) throws Exception {
    if (args.length != 3) {
      System.err.println("Usage: DataCleaning <path to NYC historic data> <path to NYC current data> <output path>");
      System.exit(-1);
    }

    Job job = Job.getInstance();
    job.setJarByClass(DataCleaning.class);
    job.setJobName("Data Cleaning");

    MultipleInputs.addInputPath(job, new Path(args[0]), TextInputFormat.class, DataCleaningNYCHistoricMapper.class);
    MultipleInputs.addInputPath(job, new Path(args[1]), TextInputFormat.class, DataCleaningNYCCurrentMapper.class);
    FileOutputFormat.setOutputPath(job, new Path(args[2]));

    job.setReducerClass(DataCleaningReducer.class);

    job.setNumReduceTasks(1);
    job.setOutputKeyClass(Text.class);
    job.setOutputValueClass(NullWritable.class);
    job.setMapOutputKeyClass(NullWritable.class);
    job.setMapOutputValueClass(Text.class);

    System.exit(job.waitForCompletion(true) ? 0 : 1);
  }
}
