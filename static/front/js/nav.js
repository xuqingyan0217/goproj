$(document).ready(function() {
    // 当鼠标悬停在导航链接上时添加 .current 类
    $("nav ul li a").hover(
        function() {
            $(this).parent().addClass("current");
        },
        function() {
            $(this).parent().removeClass("current");
        }
    );
});