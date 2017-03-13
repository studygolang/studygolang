/**
 * @license Copyright (c) 2003-2015, CKSource - Frederico Knabben. All rights reserved.
 * For licensing, see LICENSE.md or http://ckeditor.com/license
 */

/**
 * @fileOverview Rich code snippets for CKEditor.
 */

'use strict';

(function() {
  // Create a new plugin which registers a custom code highlighter
  // based on Prism.js in order to replace the one that comes
  // with the Code Snippet plugin.
  CKEDITOR.plugins.add('prism', {
    requires: 'codesnippet',

    init: function(editor) {
      var path = this.path;

      // Loading the prism.js style file.
      // Idea taken from codesnippet/plugin.js code.
      // Method is available only if wysiwygarea exists and
      // CKEditor is at least version 4.4.
      if (editor.addContentsCss) {
        editor.addContentsCss(path + 'lib/prism/prism_patched.min.css');
      }

      // Create a new instance of the highlighter.
      var prismHighlighter = new CKEDITOR.plugins.codesnippet.highlighter({
        init: function(ready) {
          // Load the Prism.js library asynchronously.
          CKEDITOR.scriptLoader.load(path + 'lib/prism/prism_patched.min.js', function() {
            // Notify the handler that the library has been loaded.
            ready();
          });
        },

        // Specify the supported languages.
        languages: {
          abap: 'ABAP',
          actionscript: 'ActionScript',
          apacheconf: 'Apache Conf',
          applescript: 'AppleScript',
          aspnet: 'ASP.NET',
          bash: 'Bash',
          basic: 'BASIC',
          c: 'C',
          coffeescript: 'CoffeeScript',
          cpp: 'C++',
          csharp: 'C#',
          css: 'CSS',
          d: 'D',
          dart: 'Dart',
          diff: 'Diff',
          docker: 'Docker',
          erlang: 'Erlang',
          fortran: 'Fortran',
          fsharp: 'F#',
          git: 'Git',
          go: 'Go',
          groovy: 'Groovy',
          haskell: 'Haskell',
          markup: 'HTML',
          http: 'HTTP',
          ini: 'INI',
          java: 'Java',
          javascript: 'JavaScript',
          lua: 'Lua',
          makefile: 'Makefile',
          markdown: 'Markdown',
          matlab: 'MATLAB',
          nginx: 'Nginx',
          objectivec: 'Objective-C',
          pascal: 'Pascal',
          perl: 'Perl',
          php: 'PHP',
          prolog: 'Prolog',
          python: 'Python',
          puppet: 'Puppet',
          r: 'R',
          ruby: 'Ruby',
          rust: 'Rust',
          sas: 'SAS',
          scala: 'Scala',
          scheme: 'Scheme',
          sql: 'SQL',
          swift: 'Swift',
          twig: 'Twig',
          vim: 'vim',
          yaml: 'YAML',
        },

        highlighter: function(code, language, callback) {
          // _self.Prism is a global namespace/object created by Prism.js.
          var _prism = _self.Prism;

          // Let the Prism.js highlight the code.
          var highlightedCode = _prism.highlight(code, _prism.languages[language], language);

          // The clever idea below is taken from the 'Line Numbers' plugin
          // of Prism. Basically, we want to count the number of newlines (\n)
          // in the highlighted code, then create the same number
          // of SPAN elements, append them to the highlighted code
          // and finally number/label them using prism.css.
          var match = highlightedCode.match(/\n(?!$)/g);
          var linesNum = match ? match.length + 1 : 1;

          var lines = new Array(linesNum + 1);
          lines = lines.join('<span></span>');

          // Create the SPAN root/wrapper, insert its child SPAN lines,
          // then append them to the highlighted code.
          var lineNumbersWrapper = '<span class="line-numbers-rows">';
          lineNumbersWrapper += lines;
          lineNumbersWrapper += '</span>';
          highlightedCode += lineNumbersWrapper;

          // Return highlighted code.
          callback(highlightedCode);
        }
      });

      // From now on, prismHighlighter will be used as a Code Snippet
      // highlighter, overwriting the default engine.
      editor.plugins.codesnippet.setHighlighter(prismHighlighter );
    }
  });
})();
