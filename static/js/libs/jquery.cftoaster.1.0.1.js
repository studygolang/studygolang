/**
 * This plugin creates a toast message
 * multiple toast messages are added to a queue and displayed one at a time
 */
(function($) {
	
	$.fn.cftoaster = function(options) 
	{
		var opts = $.extend({}, $.fn.cftoaster.options, options);
		
		return this.each(function() {
			opts.element = $(this);
			
			if(!isDestroyCommand(opts)) {
				$.cftoaster._addToQueue(opts);
			}
			
			//destroy speech bubble 
			else {
				$.cftoaster._destroy(opts);
			}
	    });
	
		/**
		 * check to see if options is a destroy command
		 * @param options, object
		 * @return isDestroyCommand, boolean
		 */
	    function isDestroyCommand(options) 
	    {
	    	var command = "";
	    	for(var i=0; i<=$.cftoaster.DESTROY_COMMAND.length; i++)
	    	{
	    		if(!options.hasOwnProperty(i)){
	    			break;
	    		}
	    		command += options[i];
	    	}
	    	return command == $.cftoaster.DESTROY_COMMAND;
	    }
	}
	
	/**
	 * customization options
	 * @param content, string, can be html
	 * @param element, DOM element to insert the message
	 * @param animationTime, int time in ms for the animation, -1 to show no animation
	 * @param showTime, int time in ms for the toast message to stay visible
	 * @param maxWidth, int maximum width of the message container in px
	 * @param backgroundColor, string, hexadecimal value of the colour, requires "#" prefix
	 * @param fontColor, string, hexadecimal value of the colour, requires "#" prefix
	 * @param bottomMargin, int, space to leave between the bottom of the toast message and the bottom of the browser window in px
	 */
	$.fn.cftoaster.options = 
	{
		content: "This is a toast message eh",
		element: "body",
		animationTime: 150,
		showTime: 3000,
		maxWidth: 250,
		backgroundColor: "#1a1a1a",
		fontColor: "#eaeaea",
		bottomMargin: 75
	}

})(jQuery);

/**
 * cftoaster specific functions
 */
jQuery.extend({
	
    cftoaster: 
    {
    	/**** constants ****/
    	NAMESPACE : "cf_toaster",
    	DESTROY_COMMAND : "destroy",
    	MAIN_CSS_CLASS : "cf_toaster",
    	
    	/**** private variables ****/
    	_queue : [],
    	
    	/**
    	 * add an item into the queue
    	 * @param options, object
    	 */
    	_addToQueue: function(options) 
    	{
    		this._queue.push(options);
    		
    		if(options.element && 
    			!this._isShowingToastMessage(options.element))
    		{
    			this._showNextInQueue(options.element);
    		}
    	},
    	
    	/**
    	 * remove all items with the specified element from queue
    	 * if no element is specified, remove all items from queue
    	 * @param element, DOM element
    	 */
    	_removeFromQueue: function(element)
    	{
    		if(element)
    		{
    			for(var i in this._queue)
    			{
    				var options = this._queue[i];
    				if( $(options.element).is(element) ){
    					this._queue.splice(i, 1);
    				}
    			}
    		}
    		else 
    		{
    			this._queue = [];
    		}
    	},
    	
    	/**
    	 * destroy the current message and remove all items from queue
    	 * if no element is specified, destroy all instances of the toast message
    	 * @param options, object
    	 */
    	_destroy: function(options)
    	{
    		var element = options && options.element ? options.element : undefined;
    		if(element){
    			$(element).find("." + this.MAIN_CSS_CLASS).remove();
    		}
    		else{
    			$("." + this.MAIN_CSS_CLASS).remove();
    		}
    		this._removeFromQueue(element);
    	},
    	
    	/**
    	 * returns if the given element currently has a toast message showing
    	 * @param element, DOM element
    	 * @return isShowing, boolean
    	 */
    	_isShowingToastMessage: function(element)
    	{
    		var isShowing = false;
    		if(element){
    			isShowing = $(element).find("." + this.MAIN_CSS_CLASS).size() > 0;
    		}
    		return isShowing;
    	},
    	
    	/**
    	 * show the next toast message in the queue
    	 * @param element, DOM element
    	 */
    	_showNextInQueue: function(element)
    	{
    		var nextItem;
    		for(var i=0; i<this._queue.length; i++)
    		{
    			var item = this._queue[i];
    			if( $(item.element).is(element) )
    			{
    				nextItem = item;
    				this._queue.splice(i, 1);
    				break;
    			}
    		}
    		
    		if(nextItem)
    		{
    			var background = $("<div/>").addClass("background").css("background", nextItem.backgroundColor);
    			var content = $("<div/>").addClass("content").html(nextItem.content)
    									 .css("width", nextItem.maxWidth + "px")
    									 .css("color", nextItem.fontColor);
    			var container = $("<div/>").addClass(this.MAIN_CSS_CLASS).hide().append(background).append(content);
    			$(element).append(container);
    			
    			//center main container
    			var marginLeft = -$(container).outerWidth()/2 + "px";
    			$(container).css("bottom", nextItem.bottomMargin + "px").css("margin-left", marginLeft);
    			
    			//animate to show and then hide
    			$(container).stop().fadeIn(nextItem.animationTime).delay(nextItem.showTime)
    						.fadeOut(nextItem.animationTime, function(){
    							$(this).remove();
    							$.cftoaster._showNextInQueue(element);
    						});
    		}
    	},
    	
    	/**
    	 * change the defaults for the plugin
    	 * any message created prior to this call retains the original default options
    	 * @param newDefaults, object
    	 */
    	setDefaults: function(newDefaults)
    	{
    		var defaults = $.extend({}, $.fn.cftoaster.options, newDefaults);
    		$.fn.cftoaster.options = defaults;
    	}
		
    }
    
});






