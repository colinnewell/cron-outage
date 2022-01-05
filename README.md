# cron-outage

Given an outage window, figure out which cron jobs would have run in that time.

    ./cron-outage --start 2022-01-03T05:30:00Z --end 2022-01-03T11:30:00Z \
        --not-before 2022-01-03T19:30:00Z cron-file

Use --start and --end to specify the outage period.

Use --not-before to specify jobs to weed out because they
will be executed regardless because they will be done soon.  Set it to an hour
or 2 after for example and that should weed out ones that run every few minutes
or hours since they will get run normally anyway.

There's no point in trying to reschedule a job with a schedule like `* * * * *`
once things are back up.
