-- clear table
DELETE FROM analytics.event;

-- clear table
DELETE FROM analytics.task;

-- 1 write event
INSERT INTO analytics.event
(task_id, occurred_at, event_type, event_user, approvers_number)
VALUES (321, '2022-02-02 16:00:00.000000 +00:00', 'create', 'author@mail.com', 3);

-- 2 create task
INSERT INTO analytics.task
(id, status, created_at, last_mail_at, total_time, approvers_number, current_approvers_number)
VALUES (321, 'created', '2022-02-02 16:00:00.000000 +00:00', null, '0 years 0 mons 0 days 0 hours 0 mins 0.0 secs', 3, 0);

-- 3  write event
INSERT INTO analytics.event
(task_id, occurred_at, event_type, event_user, approvers_number)
VALUES (321, '2022-02-02 16:00:05.000000 +00:00', 'send_mail', 'addressee1@mail.com', null);

-- 4 mail sent
UPDATE analytics.task
SET
    status = 'waiting_response',
    last_mail_at = '2022-02-02 16:00:05.000000 +00:00'
WHERE analytics.task.id = 321;

-- 5 write event
INSERT INTO analytics.event
(task_id, occurred_at, event_type, event_user, approvers_number)
VALUES (321, '2022-02-02 16:01:05.000000 +00:00', 'approve', 'addressee1@mail.com', null);

-- 6 update task - positive click
UPDATE analytics.task
SET
    status = CASE
                 WHEN current_approvers_number + 1 = approvers_number THEN 'approved'
                 ELSE 'response_received'
        END,
    total_time = task.total_time + ('2022-02-02 16:01:05.000000 +00:00'::timestamp - last_mail_at),
    current_approvers_number =  current_approvers_number + 1
WHERE analytics.task.id = 321;

-- 7 write event
INSERT INTO analytics.event
(task_id, occurred_at, event_type, event_user, approvers_number)
VALUES (321, '2022-02-02 16:01:10.000000 +00:00', 'send_mail', 'addressee2@mail.com', null);

-- 8 mail sent
UPDATE analytics.task
SET
    status = 'waiting_response',
    last_mail_at = '2022-02-02 16:01:10.000000 +00:00'
WHERE analytics.task.id = 321;

-- 9 write event
INSERT INTO analytics.event
(task_id, occurred_at, event_type, event_user, approvers_number)
VALUES (321, '2022-02-02 16:02:10.000000 +00:00', 'reject', 'addressee2@mail.com', null);

-- 10 update task - negative click
UPDATE analytics.task
SET
    status = 'rejected',
    total_time = task.total_time + ('2022-02-02 16:02:10.000000 +00:00'::timestamp - last_mail_at),
    current_approvers_number =  current_approvers_number + 1
WHERE analytics.task.id = 321;
