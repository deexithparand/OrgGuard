package orgguard

default allow = false

allow if startswith(input.target_user, "jmd")
