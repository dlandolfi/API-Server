class ApiController < ApplicationController
  def newsfeed
    render json:{message:"newsfeed"}
  end
end
