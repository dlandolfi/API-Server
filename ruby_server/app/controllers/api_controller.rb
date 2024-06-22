require 'rss'
require 'open-uri'

class ApiController < ApplicationController
  def newsfeed
    rss_url = 'https://www.coffeereview.com/feed/' # Replace with the actual RSS feed URL
    rss_content = URI.open(rss_url).read
    rss = RSS::Parser.parse(rss_content, false)

    # Get the first item from the RSS feed
    first_item = rss.items.first

    # Transform the first item as needed
    transformed_item = {
      title: first_item.title,
      link: first_item.link,
      pub_date: first_item.pubDate
    }

    render json: transformed_item
  end
end

