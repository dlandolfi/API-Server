# frozen_string_literal: true

require 'test_helper'

class ApiControllerTest < ActionDispatch::IntegrationTest
  test 'should get newsfeed' do
    get api_newsfeed_url
    assert_response :success
  end
end
