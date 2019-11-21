(function() {
'use strict';

function bindToggle(el) {
	$('.toggleButton', el).click(function() {
		if ($(this).closest(".toggle, .toggleVisible")[0] != el) {
			// Only trigger the closest toggle header.
			return;
		}

		if ($(el).is('.toggle')) {
			$(el).addClass('toggleVisible').removeClass('toggle');
		} else {
			$(el).addClass('toggle').removeClass('toggleVisible');
		}
	});
}

function bindToggles(selector) {
	$(selector).each(function(i, el) {
		bindToggle(el);
	});
}

function bindToggleLink(el, prefix) {
	$(el).click(function() {
		var href = $(el).attr('href');
		var i = href.indexOf('#'+prefix);
		if (i < 0) {
			return;
		}
		var id = '#' + prefix + href.slice(i+1+prefix.length);
		if ($(id).is('.toggle')) {
			$(id).find('.toggleButton').first().click();
		}
	});
}
function bindToggleLinks(selector, prefix) {
	$(selector).each(function(i, el) {
		bindToggleLink(el, prefix);
	});
}

function toggleHash() {
	var id = window.location.hash.substring(1);
	// Open all of the toggles for a particular hash.
	var els = $(
		document.getElementById(id),
		$('a[name]').filter(function() {
			return $(this).attr('name') == id;
		})
	);

	while (els.length) {
		for (var i = 0; i < els.length; i++) {
			var el = $(els[i]);
			if (el.is('.toggle')) {
				el.find('.toggleButton').first().click();
			}
		}
		els = el.parent();
	}
}

function personalizeInstallInstructions() {
	var prefix = '?download=';
	var s = window.location.search;
	if (s.indexOf(prefix) != 0) {
		// No 'download' query string; detect "test" instructions from User Agent.
		if (navigator.platform.indexOf('Win') != -1) {
			$('.testUnix').hide();
			$('.testWindows').show();
		} else {
			$('.testUnix').show();
			$('.testWindows').hide();
		}
		return;
	}

	var filename = s.substr(prefix.length);
	var filenameRE = /^go1\.\d+(\.\d+)?([a-z0-9]+)?\.([a-z0-9]+)(-[a-z0-9]+)?(-osx10\.[68])?\.([a-z.]+)$/;
	$('.downloadFilename').text(filename);
	$('.hideFromDownload').hide();
	var m = filenameRE.exec(filename);
	if (!m) {
		// Can't interpret file name; bail.
		return;
	}

	var os = m[3];
	var ext = m[6];
	if (ext != 'tar.gz') {
		$('#tarballInstructions').hide();
	}
	if (os != 'darwin' || ext != 'pkg') {
		$('#darwinPackageInstructions').hide();
	}
	if (os != 'windows') {
		$('#windowsInstructions').hide();
		$('.testUnix').show();
		$('.testWindows').hide();
	} else {
		if (ext != 'msi') {
			$('#windowsInstallerInstructions').hide();
		}
		if (ext != 'zip') {
			$('#windowsZipInstructions').hide();
		}
		$('.testUnix').hide();
		$('.testWindows').show();
	}

	var download = "https://storage.googleapis.com/golang/" + filename;

	var message = $('<p class="downloading">'+
		'Your download should begin shortly. '+
		'If it does not, click <a>this link</a>.</p>');
	message.find('a').attr('href', download);
	message.insertAfter('#nav');

	window.location = download;
}

function updateVersionTags() {
	var v = window.goVersion;
	if (/^go[0-9.]+$/.test(v)) {
		$(".versionTag").empty().text(v);
		$(".whereTag").hide();
	}
}

function addPermalinks() {
	function addPermalink(source, parent) {
		var id = source.attr("id");
		if (id == "" || id.indexOf("tmp_") === 0) {
			// Auto-generated permalink.
			return;
		}
		if (parent.find("> .permalink").length) {
			// Already attached.
			return;
		}
		parent.append(" ").append($("<a class='permalink'>&#xb6;</a>").attr("href", "#" + id));
	}

	$("#page .container").find("h2[id], h3[id]").each(function() {
		var el = $(this);
		addPermalink(el, el);
	});

	$("#page .container").find("dl[id]").each(function() {
		var el = $(this);
		// Add the anchor to the "dt" element.
		addPermalink(el, el.find("> dt").first());
	});
}

$(document).ready(function() {
	addPermalinks();
	bindToggles(".toggle");
	bindToggles(".toggleVisible");
	bindToggleLinks(".exampleLink", "example_");
	bindToggleLinks(".overviewLink", "");
	bindToggleLinks(".examplesLink", "");
	bindToggleLinks(".indexLink", "");
	toggleHash();
	personalizeInstallInstructions();
	updateVersionTags();
});

})();