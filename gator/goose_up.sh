export base=$(pwd)
cd /home/rrochlin/boot_dev/BlogAggregator/gator/sql/schema/
goose postgres "postgres://postgres:postgres@localhost:5432/gator" up
cd $base
