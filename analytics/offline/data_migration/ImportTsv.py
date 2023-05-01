import sys
from datetime import date
from datetime import datetime

if __name__ == "__main__":
  today = date.today()
  now = datetime.now()
  current_time = now.strftime("%H%M%S")
  with open("log-{}-{}-ImportTsv.txt".format(today, current_time), "x") as logfile: 
    table = "group4_rbda_nyu_edu:crimes"
    columns = ["x", "y", "t", "d"]
    for k, arg in enumerate(sys.argv):
      if k == 0:
        continue
      logfile.write("INFO: processing `{}`...\n".format(arg))
      with open(arg) as infile:
        line_num = 0
        for line in infile:
          # remove "\n" and replace "'" by "`" since HBase doesn't accept "'" in data
          fields = line.rstrip().replace("'", "`").split('\t')
          if len(fields) != 5:
            logfile.write("WARNNING: discard line `{}`\n".format(line))
            continue
          for i in range(1, 5):
            cmd = "put '{}', '{}', 'e:{}', '{}'".format(table, fields[0], columns[i - 1], fields[i])
            print(cmd)
          line_num += 1
          if line_num % 1000 == 0:
            logfile.write("INFO: {} lines processed\n".format(line_num))
      logfile.write("INFO: `{}` processed, {} lines in total\n".format(arg, line_num))
