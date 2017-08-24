# clnd
Cleans up after Sonarr and Radarr

Add to Sonarr or Radarr as a Custom Script.

Settings > Connect > + > Custom Script 

Name: postCleaned  
On: Download  
Path: /path/to/clnd.exe  

This will make it run after Sonarr and Radarr have finished copying the file to the final destination.
