namespace cpp seraph.meizu
namespace py seraph.meizu

struct AdLocation
{
    1: i64 position_id,
    2: i64 page_id,
    3: i64 block_id,
    4: i64 rank_id,
    5: i64 category_id,
    6: i64 tag_id,
    7: i32 position,
    8: i32 position_type,
}

struct UnitRecommend
{
    1: i64 unit_id,
    2: i64 app_id,
    3: double score,
    4: double price,
    5: optional i16 kw_type,
    6: optional string keyword,
    7: optional double ctr,
    8: optional double relevance,
    9: optional i64 id,
}

struct Idea
{
    1: i32 idea_id, // 创意ID
    2: i16 source, // 创意来源, 0表示非下载类, 1表示下载类
}

struct Req
{
    1: string imei,
    2: string session_id,
    3: list<AdLocation> ad_locations,
    4: string context_info,
    5: string abtest_parameters,
    // search请求必须同时传入query和query_type
    6: optional string query,
    7: optional string query_type,
    8: optional list<string> impr_list,
    9: optional bool enable_rank = true,
    10: optional i16 source = 1,
    11: optional i16 topk = 0,
    12: optional list<Idea> ideas, // 信息流创意预估
    13: optional string req_id,
    14: optional i64 rel_app_id = 0,
    15: optional i16 business = 1,
    255: optional bool debug = false,
}


struct Rsp
{
    1: string status = '',
    2: map<i64, list<UnitRecommend>> rank_results,
    3: string version = '',
}

struct SearchRsp
{
    1: string status = '',
    2: list<UnitRecommend> rank_results,
    3: string version = '',
}

struct BinaryRsp
{
    1: string status = '',
    2: binary rank_results,
}

struct AckReq
{
    1: string message,
    2: optional string message_type,
    3: optional bool disable_databus,
}

struct AckRsp
{
    1: string status = '',
}

struct ServerImprRsp
{
    1: i16 status,
}

service Predict
{
//    Rsp predict(1:Req req),
//    SearchRsp search(1:Req req),
    BinaryRsp binary_predict(1:Req req),
    BinaryRsp binary_search(1:Req req),
    BinaryRsp relate_predict(1:Req req),
    BinaryRsp feed_predict(1:Req req),
    ServerImprRsp upload_server_impr(1:Req req),
    oneway void upload_server_impr_oneway(1:Req req),
    AckRsp ack(1:AckReq req),
    AckRsp upload_applist(1:AckReq req),
}
