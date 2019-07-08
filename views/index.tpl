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

    /*
    @param file: 文件信息
    @param successChunkNum: 已经成功上传的 chunk 数量
    @param chunksize: 每一个 chunk 的大小
    @param chunkTotalNumber: 总共有多少个 chunk
    @param MaxConcurrentChunkNumber: 同时并发上传的 chunk 数量
    @param currentChunksStartIndex: 当前开始上传的 chunk 下标
    */
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

      /*
      从 currentChunksStartIndex 到 endChunkIndex 就是当前准备并发上传的 chunk
      */
      for (var i = currentChunksStartIndex; i < endChunkIndex; i++) {

        var start = i * chunksize; //分片的起始位置
        var end = Math.min(file.size, start + chunksize); //分片的结束位置

        var form = new FormData();
        form.append("file_name", file.name);
        form.append("file_size", file.size);
        form.append("file_chunk_total_number", chunkTotalNumber); //总片数
        form.append("chunk_index", i + 1); //当前是第几个分片
        form.append("chunk_data", file.slice(start, end));

        var ajax_upload = function (
          form,
          tryCount
        ) {
          $.ajax({
            url: "/UploadBigFile",
            type: "POST",
            data: form,
            async: true,
            processData: false,
            contentType: false,
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
                    // 当前并发上传的 chunk 全部成功之后，把 endChunkIndex 作为下一轮的 currentChunksStartIndex
                    upload_chunks(file, successChunkNum, chunksize, chunkTotalNumber, MaxConcurrentChunkNumber,
                      endChunkIndex)
                  }
                }
              } else {
                console.log("上传失败");
              }
            },
            error: function (xhr, textStatus, errorThrown) {
              console.log("ajax failed, textStatus =", textStatus, ", tryCount =", tryCount);
              if (tryCount <= this.retryLimit) {
                // 如果失败了，则 sleep 5s 之后再继续尝试。
                setTimeout(function () {
                  ajax_upload(form, tryCount+1);
                }, 5000);
                return;
              }
            }
          });
        }

        ajax_upload(form, 0);
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