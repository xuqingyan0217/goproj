/* ----

# Pio Plugin
# By: Dreamer-Paul
# Last Update: 2021.5.6

一个支持更换 Live2D 模型的 JS 插件

本代码为奇趣保罗原创，并遵守 GPL 2.0 开源协议。欢迎访问我的博客：https://paugram.com

---- */

var Paul_Pio = function (prop) {
    var that = this;

    var current = {
        idol: 0,
        menu: document.querySelector(".pio-container .pio-action"),
        canvas: document.getElementById("pio"),
        body: document.querySelector(".pio-container"),
        root: document.location.protocol +'//' + document.location.hostname +'/'
    };

    /* - 方法 */
    var modules = {
        // 更换模型
        idol: function () {
            current.idol < (prop.model.length - 1) ? current.idol++ : current.idol = 0;
            return current.idol;
        },
        // 创建内容
        create: function (tag, prop) {
            var e = document.createElement(tag);
            if(prop.class) e.className = prop.class;
            return e;
        },
        // 随机内容
        rand: function (arr) {
            return arr[Math.floor(Math.random() * arr.length + 1) - 1];
        },
        // 创建对话框方法
        render: function (text) {
            if(text.constructor === Array){
                dialog.innerText = modules.rand(text);
            }
            else if(text.constructor === String){
                dialog.innerText = text;
            }
            else{
                dialog.innerText = "输入内容出现问题了 X_X";
            }

            dialog.classList.add("active");

            clearTimeout(this.t);
            this.t = setTimeout(function () {
                dialog.classList.remove("active");
            }, 3000);
        },
        // 移除方法
        destroy: function () {
            that.initHidden();
            localStorage.setItem("posterGirl", 0);
        },
        // 是否为移动设备
        isMobile: function () {
            var ua = window.navigator.userAgent.toLowerCase();
            ua = ua.indexOf("mobile") || ua.indexOf("android") || ua.indexOf("ios");

            return window.innerWidth < 500 || ua !== -1;
        }
    };
    this.destroy = modules.destroy;

    var elements = {
        home: modules.create("span", {class: "pio-home"}),
        // skin: modules.create("span", {class: "pio-skin"}),
        // info: modules.create("span", {class: "pio-info"}),
        // 夜间模式按钮
        // night: modules.create("span", {class: "pio-night"}),
        close: modules.create("span", {class: "pio-close"}),

        show: modules.create("div", {class: "pio-show"})
    };

    var dialog = modules.create("div", {class: "pio-dialog"});
    current.body.appendChild(dialog);
    current.body.appendChild(elements.show);

    /* - 提示操作 */
    var action = {
        // 欢迎
        welcome: function () {
            if(document.referrer !== "" && document.referrer.indexOf(current.root) === -1){
                var referrer = document.createElement('a');
                referrer.href = document.referrer;
                prop.content.referer ? modules.render(prop.content.referer.replace(/%t/, "“" + referrer.hostname + "”")) : modules.render("🎉欢迎来自 “" + referrer.hostname + "” 的朋友！");
            }
            else if(prop.tips){
                var text, hour = new Date().getHours();

                if (hour > 22 || hour <= 5) {
                    text = '🌃你是夜猫子呀？这么晚还不睡觉，明天起的来嘛';
                }
                else if (hour > 5 && hour <= 8) {
                    text = '🌞早上好！';
                }
                else if (hour > 8 && hour <= 11) {
                    text = '🌞上午好！工作顺利嘛，不要久坐，多起来走动走动哦！';
                }
                else if (hour > 11 && hour <= 14) {
                    text = '🌞中午了，工作了一个上午，现在是午餐时间！';
                }
                else if (hour > 14 && hour <= 17) {
                    text = '🌞午后很容易犯困呢，今天的运动目标完成了吗？';
                }
                else if (hour > 17 && hour <= 19) {
                    text = '🌙傍晚了！窗外夕阳的景色很美丽呢，最美不过夕阳红~';
                }
                else if (hour > 19 && hour <= 21) {
                    text = '🌙晚上好，今天过得怎么样？';
                }
                else if (hour > 21 && hour <= 23) {
                    text = '🌙已经这么晚了呀，早点休息吧，晚安~';
                }
                else{
                    text = "💦小言u说：这个是无法被触发的吧，哈哈";
                }

                modules.render(text);
            }
            else{
                modules.render(prop.content.welcome || "🎉欢迎来到本站！");
            }
        },
        // 触摸
        touch: function () {
            current.canvas.onclick = function () {
                modules.render(prop.content.touch || ["😊你好呀", "🥳又是开心的一天呢~", "🤨你在干什么？", "👀达咩 达咩 达咩yo~", "😖不可以这样欺负我啦！", "🍕遇到困难不要怕！"]);
            };
        },
        // 右侧按钮
        buttons: function () {
            // 返回首页
            elements.home.onclick = function () {
                location.href = current.root;
            };
            elements.home.onmouseover = function () {
                modules.render(prop.content.home || "🏚点击回到首页");
            };
            current.menu.appendChild(elements.home);

            // 更换模型
            // elements.skin.onclick = function () {
            //     loadlive2d("pio", prop.model[modules.idol()]);
            //     prop.content.skin && prop.content.skin[1] ? modules.render(prop.content.skin[1]) : modules.render("新衣服真漂亮~");
            // };
            // elements.skin.onmouseover = function () {
            //     prop.content.skin && prop.content.skin[0] ? modules.render(prop.content.skin[0]) : modules.render("想看看我的新衣服吗？");
            // };
            // if(prop.model.length > 1) current.menu.appendChild(elements.skin);

            // 关于我
            // elements.info.onclick = function () {
            //     window.open(prop.content.link || "https://paugram.com/coding/add-poster-girl-with-plugin.html");
            // };
            // elements.info.onmouseover = function () {
            //     modules.render("想了解更多关于我的信息吗？");
            // };
            // current.menu.appendChild(elements.info);

            // 夜间模式
            if(prop.night){
                elements.night.onclick = function () {
                    typeof prop.night === "function" ? prop.night() : eval(prop.night);
                };
                elements.night.onmouseover = function () {
                    modules.render("夜间点击这里可以保护眼睛呢");
                };
                current.menu.appendChild(elements.night);
            }

            // 关闭看板娘
            elements.close.onclick = function () {
                modules.destroy();
            };
            elements.close.onmouseover = function () {
                modules.render(prop.content.close || "😟下次再见吧~Bye~");
            };
            current.menu.appendChild(elements.close);
        },
        custom: function () {
            prop.content.custom.forEach(function (t) {
                // 默认类型为 default
                if(!t.type) t.type = "default";

                // 使用更精确的选择器查询
                var elements = document.querySelectorAll(t.selector);

                if(elements.length){
                    elements.forEach(function(element) {
                        // 针对 read 类型的特殊处理
                        if(t.type === "read"){
                            element.addEventListener('mouseover', function() {
                                var text = "想阅读 %t 吗？".replace(/%t/, "“" + this.innerText.trim() + "”");
                                modules.render(text);
                            });
                        }
                        // 针对 link 类型的处理
                        else if(t.type === "link"){
                            element.addEventListener('mouseover', function() {
                                var text = "想了解一下 %t 吗？".replace(/%t/, "“" + this.innerText.trim() + "”");
                                modules.render(text);
                            });
                        }
                        // 自定义文本
                        else if(t.text){
                            element.addEventListener('mouseover', function() {
                                modules.render(t.text);
                            });
                        }
                    });
                }
            });
        }
    };

    /* - 运行 */
    var begin = {
        static: function () {
            current.body.classList.add("static");
        },
        fixed: function () {
            action.touch(); action.buttons();
        },
        draggable: function () {
            action.touch(); action.buttons();

            var body = current.body;
            body.onmousedown = function (downEvent) {
                var location = {
                    x: downEvent.clientX - this.offsetLeft,
                    y: downEvent.clientY - this.offsetTop
                };

                function move(moveEvent) {
                    body.classList.add("active");
                    body.classList.remove("right");
                    body.style.left = (moveEvent.clientX - location.x) + 'px';
                    body.style.top  = (moveEvent.clientY - location.y) + 'px';
                    body.style.bottom = "auto";
                }

                document.addEventListener("mousemove", move);
                document.addEventListener("mouseup", function () {
                    body.classList.remove("active");
                    document.removeEventListener("mousemove", move);
                });
            };
        }
    };

    // 运行
    this.init = function (onlyText) {
        if(!(prop.hidden && modules.isMobile())){
            if(!onlyText){
                action.welcome();
                loadlive2d("pio", prop.model[0]);
            }

            switch (prop.mode){
                case "static": begin.static(); break;
                case "fixed":  begin.fixed(); break;
                case "draggable": begin.draggable(); break;
            }

            if(prop.content.custom) action.custom();
        }
    };

    // 隐藏状态
    this.initHidden = function () {
        // ! 清除预设好的间距
        if(prop.mode === "draggable"){
            current.body.style.top = null;
            current.body.style.left = null;
            current.body.style.bottom = null;
        }

        current.body.classList.add("hidden");
        dialog.classList.remove("active");

        elements.show.onclick = function () {
            current.body.classList.remove("hidden");
            localStorage.setItem("posterGirl", 1);
            that.init();
        }
    }

    localStorage.getItem("posterGirl") == 0 ? this.initHidden() : this.init();
};

// 请保留版权说明
// if (window.console && window.console.log) {
//     console.log("%c 看板娘的诞生离不开Pio及Live2D %c http://zxz.ee/live2d.html ","color: #fff; margin: 1em 0; padding: 5px 0; background: #673ab7;","margin: 1em 0; padding: 5px 0; background: #efefef;");
// }