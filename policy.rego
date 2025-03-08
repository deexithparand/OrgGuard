package orgguard

# By default, deny access unless explicitly allowed
default allow = false

# Allow access only if the target user's username starts with "jmd"
allow if startswith(input.target_user, "jmd")
