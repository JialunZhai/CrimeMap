# migrate data from Hive to HBase, the schema will be changed
# install maven
  https://maven.apache.org/download.cgi
# run the following commands in this directory
# build 
  mvn package
# run
## ${OUTPUT_DIR} should not exist before run this command
  hadoop jar dataMigration.jar DataMigration ${OUTPUT_DIR} ${INPUT_FILE_1} ${INPUT_FILE_2} ... ${INPUT_FILE_N}

# import data to HBase 
  python ImportTsv.py ${INPUT_DATA_TSV} | hbase shell >/dev/null

# or 
  hbase org.apache.hadoop.hbase.mapreduce.ImportTsv -Dimporttsv.columns=HBASE_ROW_KEY,e:x,e:y,e:t,e:d group04_rbda_nyu_edu:crimes ${INPUT_DATA_TSV}