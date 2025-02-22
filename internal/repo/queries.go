package repo

var (
	listUserQ = `
	SELECT 
		COALESCE(u.userid, '00000000-0000-0000-0000-000000000000') AS userid,
		COALESCE(u.name, '') AS name,
		COALESCE(u.mailid, '') AS mailid,
		COALESCE(u.purpose, 'student') AS purpose,
		COALESCE(u.is_active, false) AS is_active,
		COALESCE(u.organization, '') AS organization,
		COALESCE(u.created_on, '1970-01-01 00:00:00') AS created_on,
		COALESCE(u.updated_on, '1970-01-01 00:00:00') AS updated_on,
		COUNT(*) OVER() AS total_count
	FROM 
		users u
	LEFT JOIN 
		user_permissions up 
	ON 
		u.userid = up.user_id
	WHERE 
		1=1`
)
