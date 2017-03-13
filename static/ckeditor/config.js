var MyEditorExtraPlugins = 'codesnippet,image2,uploadimage,notification,widget,lineutils,justify,autolink';

var MyEditorConfig = {
    title: 'Go语言中文网富文本编辑器',
    // Define the toolbar: http://docs.ckeditor.com/#!/guide/dev_toolbar
    // The standard preset from CDN which we used as a base provides more features than we need.
    // Also by default it comes with a 2-line toolbar. Here we put all buttons in a single row.
    // toolbar: [
    //     { name: 'basicstyles', items: [ 'Bold', 'Italic', 'Underline', 'RemoveFormat' ] },
    //     { name: 'paragraph', items: [ 'NumberedList', 'BulletedList', '-', 'Outdent', 'Indent', '-', 'Blockquote', '-', 'JustifyLeft', 'JustifyCenter', 'JustifyRight' ] },
    //     { name: 'clipboard', items: [ 'Undo', 'Redo' ] },
    //     { name: 'links', items: [ 'Link', 'Unlink' ] },
    //     { name: 'insert', items: [ 'CodeSnippet', 'Image' ] },
    //     { name: 'styles', items: [ 'Format' ] },
    //     { name: 'document', items: [ 'Source', 'Preview' ] },
    //     { name: 'tools', items: [ 'Maximize' ] }
    // ],
    startupFocus: true,
    // Since we define all configuration options here, let's instruct CKEditor to not load config.js which it does by default.
    // One HTTP request less will result in a faster startup time.
    // For more information check http://docs.ckeditor.com/#!/api/CKEDITOR.config-cfg-customConfig
    customConfig: '',
    // Enabling extra plugins, available in the standard-all preset: http://ckeditor.com/presets-all
    // extraPlugins: 'sourcedialog,preview,codesnippet,image2,uploadimage,notification,prism,widget,lineutils,justify,autolink',
    // Remove the default image plugin because image2, which offers captions for images, was enabled above.
    removePlugins: 'image',
    filebrowserImageUploadUrl: '/image/quick_upload?command=QuickUpload&type=Images',

    // See http://docs.ckeditor.com/#!/api/CKEDITOR.config-cfg-codeSnippet_theme
    codeSnippet_theme: 'monokai_sublime',//'ir_black',
    codeSnippet_languages: {
        go: 'Go',
        php: 'PHP',
        bash: 'Bash',
        cpp: 'C/C++',
        json: 'JSON',
        html: 'HTML',
        http: 'HTTP',
        ini: 'INI',
        java: 'Java',
        javascript: 'JavaScript',
        markdown: 'Markdown',
        nginx: 'Nginx',
        sql: 'SQL',
        yaml: 'YAML',
        armasm: 'ARM Assembly'
    },
    /*********************** File management support ***********************/
    // In order to turn on support for file uploads, CKEditor has to be configured to use some server side
    // solution with file upload/management capabilities, like for example CKFinder.
    // For more information see http://docs.ckeditor.com/#!/guide/dev_ckfinder_integration
    // Uncomment and correct these lines after you setup your local CKFinder instance.
    // filebrowserBrowseUrl: 'http://example.com/ckfinder/ckfinder.html',
    // filebrowserUploadUrl: 'http://example.com/ckfinder/core/connector/php/connector.php?command=QuickUpload&type=Files',
    /*********************** File management support ***********************/
    // Make the editing area bigger than default.
    height: 361,
    width: '98%',
    // An array of stylesheets to style the WYSIWYG area.
    // Note: it is recommended to keep your own styles in a separate file in order to make future updates painless.
    // contentsCss: [ 'https://cdn.ckeditor.com/4.6.2/standard-all/contents.css', 'mystyles.css' ],
    // Reduce the list of block elements listed in the Format dropdown to the most commonly used.
    format_tags: 'p;h1;h2;h3;h4;pre',
    // Simplify the Image and Link dialog windows. The "Advanced" tab is not needed in most cases.
    removeDialogTabs: 'image:advanced;link:advanced;link:target'
}