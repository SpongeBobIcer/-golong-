// 稍作封装，使之在引用后可以直接被使用（而不用每次都 layui.use）
function layMsg(content, options, end){
    layui.use('layer', function(){
        var layer = layui.layer;
        layer.msg(content, options, end);
    });
}
