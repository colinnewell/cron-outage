# cron-outage

Given an outage window, figure out which cron jobs would have run in that time.

    ./cron-outage --start 2022-01-03T05:30:00Z --end 2022-01-03T11:30:00Z cron-file

Note: also need to figure out which ones will run soon so they don't need to be
rescheduled.
