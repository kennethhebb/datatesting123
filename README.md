# datatesting123


This repo is a demonstration of a repository managed database, changes to the database are made through pull request to the repo.

This concept can be built out to include migrating existing data into the new deployment.

The purpose for this experiment is to develop the idea of a PR managed database that can have production and staging versions be hosted on a provider like Heroku
and then developed using review apps or similar development deployments from other providers. This way a database can be a uniform data plane across multiple API's 
and developers are not burdened with maintaining the database within their development environment, they would instead only have to manage an environment file that 
provides access to the appropriate review app deployment of the database branch. 
