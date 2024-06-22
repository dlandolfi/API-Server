require 'rss'
require 'open-uri'

class ApiController < ApplicationController
  def newsfeed
    rss_url = 'https://www.coffeereview.com/feed/' # Replace with the actual RSS feed URL
    rss_content = URI.open(rss_url).read
    rss = RSS::Parser.parse(rss_content, false)

    transformed_items = rss.items.map do |item|
      {
        title: item.title,
        link: item.link,
        pub_date: item.pubDate
      }
    end

    render json: transformed_items
  end
end

