$(function() {
    marked = SG.markSettingNoHightlight();

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
});