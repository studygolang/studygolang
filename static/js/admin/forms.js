/*
 * 	Additional function for forms.html
 *	Written by ThemePixels	
 *	http://themepixels.com/
 *
 *	Copyright (c) 2012 ThemePixels (http://themepixels.com)
 *	
 *	Built for Amanda Premium Responsive Admin Template
 *  http://themeforest.net/category/site-templates/admin-templates
 */

jQuery(document).ready(function($){
	
	///// DUAL BOX /////
	var db = jQuery('#dualselect').find('.ds_arrow .arrow');	//get arrows of dual select
	var sel1 = jQuery('#dualselect select:first-child');		//get first select element
	var sel2 = jQuery('#dualselect select:last-child');			//get second select element
	
	sel2.empty(); //empty it first from dom.
	
	db.click(function(){
		var t = (jQuery(this).hasClass('ds_prev'))? 0 : 1;	// 0 if arrow prev otherwise arrow next
		if(t) {
			sel1.find('option').each(function(){
				if(jQuery(this).is(':selected')) {
					jQuery(this).attr('selected',false);
					var op = sel2.find('option:first-child');
					sel2.append(jQuery(this));
				}
			});	
		} else {
			sel2.find('option').each(function(){
				if(jQuery(this).is(':selected')) {
					jQuery(this).attr('selected',false);
					sel1.append(jQuery(this));
				}
			});
		}
	});
	
	///// FORM VALIDATION /////
	jQuery('.stdform, .stdform_q').validate({
		submitHandler: function(form){
			if ($(form).attr('id') == 'role_form') {
				return;
			}
			formAjaxSubmit(form);
		}
	});
	/*
	jQuery("#form1").validate({
		rules: {
			firstname: "required",
			lastname: "required",
			email: {
				required: true,
				email: true,
			},
			location: "required",
			selection: "required"
		},
		messages: {
			firstname: "Please enter your first name",
			lastname: "Please enter your last name",
			email: "Please enter a valid email address",
			location: "Please enter your location"
		}
	});
	*/

	// 表单ajax提交
	$("form[action-type=ajax-submit]").on('submit', function(event){
		event.preventDefault();
		formAjaxSubmit(this);
	});

	// 异步提交表单
	function formAjaxSubmit(form)
	{
		$('#loaders').show();

		that = form;

		var url = $(form).attr('action'),
			data = $(form).serialize();

		if (data) {
			data += '&';
		}
		data += 'format=json&submit=1';
		
		$.ajax({
			"url": url,
			"type": "post",
			"data" : data,
			"dataType" : "json",
			"error" : function (jqXHR, textStatus, errorThrown) {
				$('#loaders').hide();
				var errMsg = errorThrown == 'Forbidden' ? "亲，没权限呢!" : "亲，服务器忙!"; jAlert(errMsg, "提示");
			},
			"success" : function (data) {
				$('#loaders').hide();
				if (data.ok) {
					jAlert("操作成功", "信息");
				} else {
					jAlert(data.error, "出错");
					return;
				}
				// $('#tooltip').text("操作成功！");
				if (typeof formSuccCallback !== "undefined") {
					formSuccCallback(data);
				} else {
					that.reset();
				}
			}
		});
	}
	
	///// TAG INPUT /////
	
	// jQuery('#tags').tagsInput();

	
	///// SPINNER /////
	
	// jQuery("#spinner").spinner({min: 0, max: 100, increment: 2});
	
	
	///// CHARACTER COUNTER /////

	/*
	jQuery("#textarea2").charCount({
		allowed: 120,		
		warning: 20,
		counterText: 'Characters left: '	
	});
	*/
	
	///// SELECT WITH SEARCH /////
	// jQuery(".chzn-select").chosen();
	
});