<!DOCTYPE html>

<html>

<head>
  <title>HTTP Upload Big File Demo</title>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
</head>

<body>
  <input type="file" id="file" />
  <button id="upload">上传</button>
  <span id="output" style="font-size: 12px">上传进度</span>

  <script src="https://apps.bdimg.com/libs/jquery/2.1.4/jquery.min.js">
  </script>

  <script>
    $("#upload").click(function () {
      upload();
    });

    function upload_chunks(
      file,
      successChunkNum,
      chunksize,
      chunkTotalNumber,
      MaxConcurrentChunkNumber,
      currentChunksStartIndex
    ) {
      var currentSuccessChunkNumber = 0;
      var startChunkIndex = currentChunksStartIndex;
      var endChunkIndex = currentChunksStartIndex + MaxConcurrentChunkNumber;
      if (endChunkIndex > chunkTotalNumber) {
        endChunkIndex = chunkTotalNumber;
      }

      for (var i = currentChunksStartIndex; i < endChunkIndex; i++) {

        var start = i * chunksize; //分片的起始位置
        var end = Math.min(file.size, start + chunksize); //分片的结束位置

        var form = new FormData();
        form.append("file_name", file.name);
        form.append("file_size", file.size);
        form.append("file_chunk_total_number", chunkTotalNumber); //总片数
        form.append("chunk_index", i + 1); //当前是第几个分片
        form.append("chunk_data", file.slice(start, end));

        $.ajax({
          url: "/UploadBigFile",
          type: "POST",
          data: form,
          async: true,
          processData: false,
          contentType: false,
          tryCount: 0,
          retryLimit: 300000,
          success: function (data) {
            console.log("data = ", data);
            if (data.errno === 200) {
              ++successChunkNum;
              ++currentSuccessChunkNumber;
              $("#output").text(successChunkNum + " / " + chunkTotalNumber);
              if (successChunkNum === chunkTotalNumber) {
                console.log("全部上传完成");
              } else {
                if (currentSuccessChunkNumber == MaxConcurrentChunkNumber) {
                  upload_chunks(file, successChunkNum, chunksize, chunkTotalNumber, MaxConcurrentChunkNumber,
                    endChunkIndex)
                }
              }
            } else {
              console.log("上传失败");
            }
          },
          error: function (xhr, textStatus, errorThrown) {
            this.tryCount++;
            console.log("textStatus =", textStatus, ", tryCount =", this.tryCount);
            if (this.tryCount <= this.retryLimit) {
              setTimeout(function() {
                $.ajax(this);
              }, 2000);
              // $.ajax(this);
              return;
            }
          }
        });
      }
    }

    function upload() {
      var file = $("#file")[0].files[0];
      if (file == undefined) {
        console.log("请先选择文件");
        return false;
      }

      var successChunkNum = 0;
      var chunksize = 100 * 1024 * 1024; //以 100M 为一个分片
      var chunkTotalNumber = Math.ceil(file.size / chunksize); //总片数
      upload_chunks(file, successChunkNum, chunksize, chunkTotalNumber, 2, 0)
    }
  </script>
</body>

</html>