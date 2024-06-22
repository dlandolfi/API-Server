# frozen_string_literal: true

require 'rss'
require 'nokogiri'
require 'open-uri'
require 'date'

class ApiController < ApplicationController
  def newsfeed
    rss_urls = ['https://www.coffeereview.com/feed/',
                'https://concretewaves.com/longboards/feed/',
                'https://www.nomadicmatt.com/feed/']

    first_items = rss_urls.map do |url|
      rss_content = URI.open(url).read
      rss = RSS::Parser.parse(rss_content, false)
      rss.items[0] # Only using the top item of each feed
    end

    # Transform the first items
    transformed_items = first_items.map do |item|
      # This is used to strip html from the description
      stripped_description = Nokogiri::HTML(item.description).text
      {
        title: item.title,
        description: "#{stripped_description[0, 50]}...",
        link: item.link,
        pub_date: item.pubDate.strftime('%B %d, %Y')
      }
    end

    render json: transformed_items
  end
end
