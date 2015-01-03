#
# Cookbook Name:: filesync
# Recipe:: default
#
# Copyright (c) 2014 Lexmark International Technology S.A.  All rights reserved.
#
# All rights reserved - Do Not Redistribute

# prerequisites
include_recipe 'annoyances'
include_recipe 'build-essential'

package('ruby2.0').run_action(:install)
package('ruby2.0-dev').run_action(:install)
package('ruby1.9.1-dev').run_action(:install)
package('libruby2.0').run_action(:install)

execute 'symlink ruby2.0' do
  command 'sudo ln -sf /usr/bin/ruby2.0 /usr/bin/ruby'
end

execute 'symlink rubygems2.0' do
  command 'sudo ln -sf /usr/bin/gem2.0 /usr/bin/gem'
end

execute 'symlink IRB' do
 command 'sudo ln -sf /usr/bin/irb2.0 /usr/bin/irb'
end

execute 'symlink RDoc' do
  command 'sudo ln -sf /usr/bin/rdoc2.0 /usr/bin/rdoc'
end

execute 'symlink ERB' do
  command 'sudo ln -sf /usr/bin/erb2.0 /usr/bin/erb'
end

package('libcurl4-openssl-dev').run_action(:install)
package('libxml2-dev').run_action(:install)
package('libxslt1-dev').run_action(:install)
package('libpq-dev').run_action(:install)

node.set['postgresql']['password']['postgres'] = 'password'

include_recipe 'postgresql::server'
include_recipe 'postgresql::contrib'
include_recipe 'git'
include_recipe 'filesync::javascript_runtime'

execute 'gem install bundler'
execute 'gem install unicorn'

package('redis-server').run_action(:install)

