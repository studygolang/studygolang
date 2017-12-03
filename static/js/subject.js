$(function() {
    marked = SG.markSetting();
    SG.Subjects = function(){}
    SG.Subjects.prototype = new SG.Publisher();

    $('.desc').html(marked($('.desc').text().trim()));

    $('.noavatar').each(function() {
        var author = $(this).data('author');
        var hash = md5(author+"");
        var data = new Identicon(hash, {format: 'svg', size: 32}).toString();
        var imgData = 'data:image/svg+xml;base64,' + data;
        $(this).attr('src', imgData);
    });

    var followed = $('#follow').data('follow');

    $('#follow').on('click', function() {
        var that = this;
        $.post('/subject/follow', {sid: $(this).data('sid')}, function(result) {
            if (result.ok) {
                if (followed) {
                    followed = false;
                    $(that).removeClass('btn-followed').addClass('btn-success');
                    $(that).html('<i class="fa fa-plus" aria-hidden="true"></i> 关注');
                } else {
                    followed = true;
                    $(that).removeClass('btn-success').addClass('btn-followed');
                    $(that).html('<i class="fa fa-check" aria-hidden="true"></i> 已关注');
                }
            }
        });
    });

    $('#follow').on('mouseenter', function() {
        if (!followed) {
            return;
        }
        $(this).html('<i class="fa fa-remove" aria-hidden="true"></i> 取消关注');
    })

    $('#follow').on('mouseleave', function() {
        if (!followed) {
            return;
        }
        $(this).html('<i class="fa fa-check" aria-hidden="true"></i> 已关注');
    })

    $('#contribute').on('click', function() {
        var sid = $('#follow').data('sid');
        $.getJSON('/subject/my_articles?sid='+sid, function(result) {
            if (result.ok) {
                var articles = result.data.articles;
                fillArticles(articles);

                $('body').addClass('modal-open');
                $('.contribute-modal').fadeIn();
            }
        });
    });

    $('.contribute-modal .close').on('click', function() {
        $('body').removeClass('modal-open');
        $('.contribute-modal').fadeOut();
    })

    var noteListHtml = '';
    $('.contribute-modal .search-btn').on('click', function() {
        var kw = $('.contribute-modal .search-input').val();
        if (kw == "") {
            $('#contribute-note-list').html(noteListHtml);
            return;
        }

        noteListHtml = $('#contribute-note-list').html();
        $('#contribute-note-list').html('');
        var placeholder = $('.contribute-modal .modal-notes-placeholder');
        placeholder.show();
        
        var sid = $('#follow').data('sid');
        $.getJSON('/subject/my_articles?kw='+encodeURIComponent(kw)+'&sid='+sid, function(result) {
            placeholder.hide();
            
            if (result.ok) {
                var articles = result.data.articles;
                if (articles.length == 0) {
                    $('#contribute-note-list').html('<div class="default">未找到相关文章</div>');
                } else {
                    fillArticles(articles);
                }
            } else {
                $('#contribute-note-list').html('<div class="default">'+result.msg+'</div>');
            }
        })
    })

    $('.contribute-modal .search-input').on('change', function() {
        if ($(this).val() == '') {
            $('#contribute-note-list').html(noteListHtml);
        }
    });

    $(document).keypress(function(evt){
        if (evt.which == 10 || evt.which == 13) {
            $('.contribute-modal .search-btn').click();
        }
    });

    $('.contribute-modal').on('click', '.action-btn', function() {
        var $articleDiv = $(this).parent(),
            articleId = $articleDiv.data('id'),
            sid = $('#follow').data('sid');
        
        var that = this;

        if ($(this).hasClass('push')) {
            $.post('/subject/contribute', {sid: sid, article_id: articleId}, function(result) {
                if (result.ok) {
                    $articleDiv.children('.note-name').addClass('has-add');
                    $(that).removeClass('push').addClass('remove').
                        before('<span class="status has-add">已加入</span>').text('移除');
                } else {
                    alert(result.error);
                }
            });
        } else {
            $.post('/subject/remove_contribute', {sid: sid, article_id: articleId}, function(result) {
                if (result.ok) {
                    $articleDiv.children('.note-name').removeClass('has-add');
                    $(that).removeClass('remove').addClass('push').text('投稿');
                    $articleDiv.children('.status').remove();
                } else {
                    alert(result.error);
                }
            });
        }
    });

    // 发布主题
    $('#submit').on('click', function(evt){
        evt.preventDefault();
        var validator = $('.validate-form').validate();
        if (!validator.form()) {
            return false;
        }

        var subjects = new SG.Subjects();
        subjects.publish(this, function(data) {
            purgeComposeDraft(uid, 'subject');

            setTimeout(function(){
                if (data.sid) {
                    window.location.href = '/subject/'+data.sid;
                } else {
                    window.location.href = '/subjects';
                }
            }, 1000);
        });
    });

    function fillArticles(articles) {
        var listHtml = '';
        for(var i in articles) {
            listHtml += '<li>'+
                '<div class="article-div" data-id="'+articles[i].id+'">';
            
            if (articles[i].had_add) {
                listHtml += '<div class="note-name has-add">'+articles[i].title+'</div>'+
                    '<span class="status has-add">已加入</span>'+
                    '<a class="action-btn remove">移除</a>';
            } else {
                
                listHtml += '<div class="note-name">'+articles[i].title+'</div>'+
                    '<a class="action-btn push">投稿</a>';
            }

            listHtml += '</div></li>';
        }
        $('#contribute-note-list').html(listHtml);
    }

    if (typeof plupload != "undefined") {
        // 实例化一个plupload上传对象
        var uploader = new plupload.Uploader({
            browse_button : 'upload-img', // 触发文件选择对话框的按钮，为那个元素id
            url : '/image/upload', // 服务器端的上传页面地址
            filters: {
                mime_types : [ //只允许上传图片
                    { title : "图片文件", extensions : "jpg,png" }
                ],
                max_file_size : '500k', // 最大只能上传 500kb 的文件
                prevent_duplicates : true // 不允许选取重复文件
            },
            multi_selection: false,
            file_data_name: 'img',
            resize: {
                width: 80
            }
        });

        // 在实例对象上调用init()方法进行初始化
        uploader.init();

        uploader.bind('FilesAdded',function(uploader, files){
            // 调用实例对象的start()
            uploader.start();
        });
        uploader.bind('UploadProgress',function(uploader,file){
            // 上传进度
        });
        uploader.bind('FileUploaded',function(uploader,file,responseObject){
            if (responseObject.status == 200) {
                var data = $.parseJSON(responseObject.response);
                if (data.ok) {
                    var url = data.data.url;
                    var path = data.data.uri;
                    var $img = $('#img-preview').find('img');
                    $img.attr('src', url);
                    $img.attr('alt', file.name);
                    $('#img-preview').show();

                    $('#cover').val(url);

                } else {
                    comTip("上传失败："+data.error);
                }
            } else {
                comTip("上传失败：HTTP状态码："+responseObject.status);
            }
        });
        uploader.bind('Error',function(uploader,errObject){
            comTip("上传出错了："+errObject.message);
        });
    }
});