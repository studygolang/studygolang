$(function(){
    $('#markdown-content').on('keydown', function(e) {
        if (e.keyCode == 9) {
            e.preventDefault();
            var indent = "\t";
            var start = this.selectionStart;
            var end = this.selectionEnd;
            var selected = window.getSelection().toString();
            selected = indent + selected.replace(/\n/g, '\n' + indent);
            this.value = this.value.substring(0, start) + selected
                    + this.value.substring(end);
            this.setSelectionRange(start + indent.length, start
                    + selected.length);
        }
    });

    $('#markdown-content').on('input propertychange', function() {
        var markdownString = $(this).val();

        // 配置 marked 语法高亮
        marked = SG.markSettingNoHightlight();

        var contentHtml = marked(markdownString);
        contentHtml = SG.replaceCodeChar(contentHtml);
        
        $('#content-preview').html(contentHtml);
        Prism.highlightAll();

        // emoji 表情解析
        emojify.run($('#content-preview').get(0));
    });

    $('#markdown-content').pasteUploadImage('/image/paste_upload');
});