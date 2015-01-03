# Cookbook Name:: filesync
# Recipe:: javascript_runtime
#
# Copyright (c) 2014 Lexmark International Technology S.A.  All rights reserved.
# All rights reserved - Do Not Redistribute
#

package 'nodejs'
package 'npm'

link '/usr/bin/node' do
  to '/usr/bin/nodejs'
end
