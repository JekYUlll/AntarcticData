$(function(){

    // 响应
    let mate_width = `<meta name="viewport" content="width=1600,target-densitydpi=device-dpi">`
    let max_client_width = $(document.body).width();
    console.log(max_client_width)
    if (max_client_width < 1200){
        $('head').append(mate_width)
    }

    // 文章
    var news_chose = $('.news-rg-chose li a')
    function detale_news_active(){
        news_chose.removeClass('news-active')
    }  
    function news_content_none(){
        $('.dtxw-content').css('display','none')
        $('.djgz-content').css('display','none')
        $('.gsgg-content').css('display','none')
        $('.mtsm-content').css('display','none')
    }
    news_chose.on('click',function(){
        detale_news_active()
        $(this).addClass('news-active')
        let attr_name = $(this).attr('data-name')
        news_content_none()
        if (attr_name == 'djgz'){
            $('.djgz-content').css('display','block')
        }else if(attr_name == 'gsgg'){
            $('.gsgg-content').css('display','block')
        }else if(attr_name == 'mtsm'){
            $('.mtsm-content').css('display','block')
        }else {
            $('.dtxw-content').css('display','block')
        }
    })

    // 当前时间
    let now_time = new Date()
    var year = now_time.getFullYear();
    var mon = now_time.getMonth() + 1;
    var date = now_time.getDate();
    var h = now_time.getHours();
    var m = now_time.getMinutes();
    // $('.ccz').text(year+'.'+mon+'.'+date+' '+h+':'+m)

    // 实时数据
    let sssj_option = $('.sssj-option')
    function detale_sssjOption_active(){
        sssj_option.removeClass('sssj-active')
    }
    function sssj_rg_none(){
        $('.sssj-rg').css('display','none')
    }

    sssj_option.on('click',function(){
        detale_sssjOption_active();
        $(this).addClass('sssj-active');
        let sssj_name = $(this).attr('data-name');
        console.log(sssj_name)
        sssj_rg_none();
        if (sssj_name == 'ccz'){
            $('.czc').css('display','block')
        }else if(sssj_name == 'zsz'){
            $('.zsz').css('display','block')
        }else if(sssj_name == 'klz'){
            $('.klz').css('display','block')
        }else if(sssj_name == 'hhz'){
            $('.hhz').css('display','block')
        }else if(sssj_name == 'tsz'){
            $('.tsz').css('display','block')
        }else {
            $('.qlz').css('display','block')
        }

    })


    // 报告图集
    var swiper = new Swiper('#bgtj-swiper', {
        slidesPerView: 4,
        spaceBetween: 9,
        autoplay: {
            delay:3500,
            disableOnInteraction: false,
        },
        speed:800,
        breakpoints: {
            1650: {
                slidesPerView: 3,
            }
        }
    });


    // 新闻轮播
    var news_swiper = new Swiper("#news-swiper", {
        spaceBetween: 30,
        pagination: {
          el: "#news-page",
          clickable: true,
        },
        autoplay:{
            delay:3500,
            disableOnInteraction: false,
        },
        speed:800,
        effect: "fade"
      });
      
      
      
      
    // 实时风向
    // index索引值，value角度
    function now_wind_direction(index,value){
        $('.sssj-rg').eq(index).find('.fx-box p').css({
            transform:'rotate('+value+'deg)'
        })
    }
    
    // 实时数据风扇 1125  历史风力级别0-17
    // index索引值，grade风力
    function now_wind(index,grade){
        if(grade == 0){return}
        var time = grade>17?0:60-grade*3
        var speed = 1
        setInterval(() => {
            $('.sssj-rg').eq(index).find('.sssj-fs-content>div>img').css({
                transform:'rotate('+speed+'deg)',
                transition: 'all '+time+'ms',
            })
            speed = speed + grade*2
            if(speed > 3600 + grade){
                speed = grade
                $('.sssj-rg').eq(index).find('.sssj-fs-content>div>img').css({
                    transform:'rotate('+speed+'deg)',
                    transition: 'all 0ms',
                })  
            }
        },time);
    }



    // 这是测试代码，数据的设置还请不要使用点击事件来添加。否则会导致定时器不断叠加，转速越点击越快
    $('.sssj-option').each(function(index){
        now_wind_direction(index,$('.sssj-rg').eq(index).find('.sssj-fx-content>div>span').text()-90)
        now_wind(index,$('.sssj-rg').eq(index).find('.sssj-fs-content>p').attr('data-fs'))
    })

})