#!/bin/sh

# args
ALL_DATA=0

# parse cmdline args
for arg in "$@"; do
  case $arg in
    -a|--all_data)
      ALL_DATA=1
      shift
      ;;
    -h|--help)
      echo "Info: run `./start.sh -a` to run MapReduce with all data from source websites"
      echo "Info: run `./start.sh` to run MapReduce with sample data from NYC_historic.tsv and NYC_current.tsv"
      shift
      ;;
    -*|--*)
      echo "Error: Unknown option $arg"
      exit 1
      ;;
    *)
      ;;
  esac
done

# compile
echo "Info: compiling..."
javac -classpath `hadoop classpath` *.java
if [ $? -ne 0 ]; then
  echo "Error: failed to compile"
  exit 1
fi
# create jar file
echo "Info: creating jar file..."
jar cvf dataCleaning.jar *.class
if [ $? -ne 0 ]; then
  echo "Error: failed to create the jar file"
  exit 1
fi

# put input files to HDFS
hadoop fs -rm -r "data_cleaning"
hadoop fs -mkdir data_cleaning

if [ $ALL_DATA -ne 0 ]; then 
  echo "Info: downloading data from source websites..."
  rm -rf NYC_historic_full.tsv && rm -rf NYC_current_full.tsv
  wget -O NYC_historic_full.tsv https://data.cityofnewyork.us/api/views/qgea-i56i/rows.tsv?accessType=DOWNLOAD
  if [ $? -ne 0 ]; then
    echo "Error: failed to download historic data"
    exit 1
  fi
  wget -O NYC_current_full.tsv https://data.cityofnewyork.us/api/views/5uac-w243/rows.tsv?accessType=DOWNLOAD
  if [ $? -ne 0 ]; then
    echo "Error: failed to download current data"
    exit 1
  fi
  echo "Info: uploading data to HDFS..."
  hadoop fs -put NYC_historic_full.tsv data_cleaning
  if [ $? -ne 0 ]; then
    echo "Error: failed to put NYC_historic_full.tsv to HDFS"
    exit 1
  fi
  hadoop fs -put NYC_current_full.tsv data_cleaning
  if [ $? -ne 0 ]; then
    echo "Error: failed to put NYC_current_full.tsv to HDFS"
    exit 1
  fi
else 
  echo "Info: uploading data to HDFS..."
  hadoop fs -put NYC_historic.tsv data_cleaning
  if [ $? -ne 0 ]; then
    echo "Error: failed to put NYC_historic.tsv to HDFS"
    exit 1
  fi
  hadoop fs -put NYC_current.tsv data_cleaning
  if [ $? -ne 0 ]; then
    echo "Error: failed to put NYC_current.tsv to HDFS"
    exit 1
  fi
fi

# run
echo "Info: running MapReduce pragram..."
if [ $ALL_DATA -ne 0 ]; then 
  hadoop jar dataCleaning.jar DataCleaning data_cleaning/NYC_historic_full.tsv data_cleaning/NYC_current_full.tsv data_cleaning/output
else 
  hadoop jar dataCleaning.jar DataCleaning data_cleaning/NYC_historic.tsv data_cleaning/NYC_current.tsv data_cleaning/output
fi

if [ $? -ne 0 ]; then
  echo "Error: failed to run MapReduce program"
  exit 1
fi

echo "Info: output files are in data_cleaning/output dir on HDFS"