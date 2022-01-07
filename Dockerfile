FROM centos
COPY images/ /images/
COPY templates/ /templates/
COPY css/ /css/
COPY levform .
EXPOSE 8080
CMD ["./levform"]
