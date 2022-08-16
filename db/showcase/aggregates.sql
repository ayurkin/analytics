-- number of approved tasks
SELECT count(id) FROM analytics.task
WHERE analytics.task.status = 'approved';

-- number of rejected tasks
SELECT count(id) FROM analytics.task
WHERE analytics.task.status = 'rejected';

-- task total response time
SELECT total_time AS total_response_time FROM analytics.task
WHERE id = 321;