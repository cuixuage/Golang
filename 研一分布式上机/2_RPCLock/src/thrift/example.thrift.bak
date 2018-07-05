namespace py example

struct Rsp {
    1: i64    cliID          //lock的client id
    2: string   operator       //get/relase operator
    3: optional string buffer
}

struct Req {
    1: i64    cliID
    2: string   operator
}

service GetLock {
    Rsp do_lock(1:Req req),
    Rsp un_lock(1:Req req),
    Rsp client_states(),

}