import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.hbase.HBaseConfiguration;
import org.apache.hadoop.hbase.KeyValue;
import org.apache.hadoop.hbase.MasterNotRunningException;
import org.apache.hadoop.hbase.ZooKeeperConnectionException;
import org.apache.hadoop.hbase.client.HBaseAdmin;
import org.apache.hadoop.hbase.client.HTable;
import org.apache.hadoop.hbase.io.ImmutableBytesWritable;
import org.apache.hadoop.hbase.mapreduce.HFileOutputFormat;
import org.apache.hadoop.hbase.mapreduce.LoadIncrementalHFiles;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;

public class DataImporting {

  static Configuration hbaseconfiguration = null;
  static Configuration conf = new Configuration();
  static HBaseAdmin hbaseAdmin;

  public static void connectHBase(String zkquorum) {
    System.out.println("Initializing Connection with Hbase");
    conf.set("hbase.zookeeper.property.clientPort", "2181");
    conf.set("hbase.zookeeper.quorum", zkquorum);

    hbaseConfiguration = HBaseConfiguration.create(conf);

    try {
      hbaseAdmin = new HBaseAdmin(hbaseConfiguration);
      System.out.println("HBase connected");
    } catch (MasterNotRunningException e) {
      System.out.println("HBase Master Exception: " + e);
    } catch (ZooKeeperConnectionException e) {
      System.out.println("Zookeeper Exception: " + e);
    }
  }

  public static void main(String[] args) throws Exception {
    if (args.length != 4) {
      System.err.println("Usage: DataImporting <input path> <output path> <zkquorum> <table name>");
      System.exit(-1);
    }
    String inputPath = args[0];
    String outputPath = args[1];
    String zkquroum = args[2];
    String tableName = args[3];

    connectHBase(zkquroum);
    conf.set("hbase.table.name", tableName);

    Job job = new Job(conf);
    job.setJarByClass(DataImporting.class);
    job.setJobName("Data Importing");

    FileInputFormat.addInputPath(job, new Path(inputPath));
    FileOutputFormat.setOutputPath(job, new Path(outputPath));

    job.setMapperClass(DataImportingMapper.class);
    job.setMapOutputKeyClass(ImmutableBytesWritable.class);
    job.setMapOutputValueClass(KeyValue.class);

    job.setNumReduceTasks(0);

    HTable htable = new HTable(conf, tableName);
    HFileOutputFormat.configureIncrementalLoad(job, htable);

    if (job.waitForCompletion(true) == 1) {
      System.exit(1);
    }

    // Importing the generated HFiles into a HBase table
    LoadIncrementalHFiles loader = new LoadIncrementalHFiles(conf);
    loader.doBulkLoad(new Path(outputPath, htable));
  }
}
