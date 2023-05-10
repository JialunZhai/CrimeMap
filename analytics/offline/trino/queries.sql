create external table crime (longitude double, latitude double, datetime bigint, description string) PARTITIONED BY (city STRING) row format delimited fields terminated by '\t' location '/user/rhd9863_nyu_edu/project/cleaned_data';

ALTER TABLE crime ADD PARTITION (city='Chicago') LOCATION '/user/rhd9863_nyu_edu/project/cleaned_data/chicago_crime';
ALTER TABLE crime ADD PARTITION (city='New York') LOCATION '/user/rhd9863_nyu_edu/project/cleaned_data/ny_crime';
ALTER TABLE crime ADD PARTITION (city='Los Angeles') LOCATION '/user/rhd9863_nyu_edu/project/cleaned_data/la_crime';
ALTER TABLE crime ADD PARTITION (city='San Francisco') LOCATION '/user/rhd9863_nyu_edu/project/cleaned_data/sf_crime';
ALTER TABLE crime ADD PARTITION (city='Austin') LOCATION '/user/rhd9863_nyu_edu/project/cleaned_data/austin_crime';
ALTER TABLE crime ADD PARTITION (city='Seattle') LOCATION '/user/rhd9863_nyu_edu/project/cleaned_data/seattle_crime';