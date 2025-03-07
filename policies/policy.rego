package orgguard

default allow = false

allow {
    startswith(input.targetUser, "jmd")
}
