window.onload = function() {

  var socket = new WebSocket("ws://golang-game-nathael.c9users.io:8082/"); //ws://10.38.162.210:8081 -- Johnathan //ws://golang-game-nathael.c9users.io:8082/ -Nathael
      socket.onopen = function(){
          console.log('Listening');
          };
      socket.onmessage = function(evt){
          var obj = JSON.parse(evt.data);
          parseKind(obj);
      };

  function socketGo(){
        var name = document.getElementById('name').value;
          var x = document.getElementById('x').value;
          var y = document.getElementById('y').value;
          var add = {
              name:name,
              body : [{x : Number(x), y :Number(y)}]
          };
            var JSON_add =  JSON.stringify(add);
          socket.send(JSON_add);
  }

  function parseKind(obj){
    /*  if(obj.kind == "restart"){
          refreshAllWs();
      }*/
      if(obj.kind == "update"){
          drawMap();
          drawSnake(obj.snakes);
          drawApple(obj.apples);
      }
      if(obj.kind == "init"){
          checkSlot(obj.players_slot);
          checkStateGame(obj.state_game);
          size =700;
          nbCellbyLine = obj.map_size;
          cellSize = size/nbCellbyLine;
      }
      if(obj.kind == "won"){
          endGame(obj.player);
      }
  }
  function checkStateGame(StateGame){
      if(StateGame == "playing"){
          document.querySelector("#start").disabled=true;
          document.querySelector("#player1").disabled=true;
          document.querySelector("#player2").disabled=true;
      }
  }



  var canvas = document.getElementById('canvas');
  if(!canvas)
  {
      alert("Impossible de récupérer le canvas");
      return;
  }

  var context = canvas.getContext('2d');
  if(!context)
  {
      alert("Impossible de récupérer le context du canvas");
      return;
  }

  function drawMap() {
      for(var x = 0 ; x <= nbCellbyLine; x++){
          for (var y = 0; y <= nbCellbyLine; y++){
              if (x%2 == 0){
                  if(y%2 == 0){
                      context.beginPath();
                      context.rect(x*cellSize,y*cellSize,cellSize,cellSize);
                      context.fillStyle = 'white';
                         context.lineWidth = 0.5;
                      context.fill();
                  }
                  else{
                      context.beginPath();
                      context.rect(x*cellSize,y*cellSize,cellSize,cellSize);
                      context.fillStyle = 'white';
                       context.lineWidth = 0.5;
                      context.fill();
                  }
              }
              else{
                  if(y%2 != 0) {
                      context.beginPath();
                      context.rect(x * cellSize, y * cellSize, cellSize, cellSize);
                      context.fillStyle = 'white';
                       context.lineWidth = 0.5;
                      context.fill();
                  }
                  else{
                      context.beginPath();
                      context.rect(x*cellSize,y*cellSize,cellSize,cellSize);
                      context.fillStyle = 'white';
                       context.lineWidth = 0.5;
                      context.fill();
                  }
              }
              context.stroke();
          }
      }
  }

/*  function refreshAllWs () {
      window.location.reload();
  }*/

  function drawSnake(S,StateGame){
      for(var i in S){
          for(var j=0; j<S[i].body.length;j++){
              if( S[i].state == "alive"){
                  context.beginPath();
                  context.rect(S[i].body[j].x*cellSize,S[i].body[j].y*cellSize,cellSize,cellSize);
                  context.fillStyle = S[i].color;
                  context.fill();
              }
              else if(StateGame == "playing"){
                  document.getElementById('dead_message').innerHTML = S[i].name + " is " + S[i].state;
              }
          }
      }
  }

  function drawApple(A) {
      for(var i in A){
          context.beginPath();
          context.rect(A[i].x*cellSize,A[i].y*cellSize,cellSize,cellSize);
          context.fillStyle = 'yellow';
          context.fill();
      }
  }


  var keySnake = {
      kind : 'move',
      key : ''
  };

  document.body.onkeydown = function(e) {
      if (e.keyCode == 90 || e.keyCode == 38) {
          //moveSnake(0, -1);
          keySnake.key = "up";
          socket.send(JSON.stringify(keySnake));
      }
      if (e.keyCode == 81 || e.keyCode == 37) {
          keySnake.key = "left";
          //moveSnake(-1, 0);
          socket.send(JSON.stringify(keySnake));
      }
      if (e.keyCode == 83 || e.keyCode == 40) {
          keySnake.key = "down";
         //moveSnake(0, 1);
          socket.send(JSON.stringify(keySnake));
      }
      if (e.keyCode == 68 || e.keyCode == 39) {
          keySnake.key = "right";
          //moveSnake(1, 0);
          socket.send(JSON.stringify(keySnake));
      }
      //collision();
      //drawMap();
      //drawApple();
      //drawSnake();
  };

  function killYourself(){
      for(var i=1; i<S.length;i++){
          if(S[0].x==S[i].x && S[0].y == S[i].y){
              return true;
          }
      }
  }

  function wallCollision(){
      if(S[0].x >19 || S[0].y >19 || S[0].x < 0 || S[0].y < 0){
          return true;
      }
  }

function endGame(winner){
  document.querySelector('#end_message').innerHTML = winner + " win";
}


  var addPlayer = {
      kind : 'connect',
      slot : '',
      name : '',
      color: '',
  };

  var start = {
      kind : 'start',
  };

  document.querySelector("#start").onclick = function(){
      socket.send(JSON.stringify(start));
  };

  document.querySelector("#player1").onclick = function(){
      addPlayer.slot = 1;
      addPlayer.name = 'Hector';
      addPlayer.color = 'black';
      socket.send(JSON.stringify(addPlayer));
      document.querySelector("#start").disabled = false;
      document.querySelector("#player2").style.display = 'none';
      document.querySelector("#player3").style.display ='none';
      document.querySelector("#player4").style.display = 'none';
      document.querySelector("#player1").disabled = true;
      document.querySelector("#spectator").disabled = true;
  };

  document.querySelector("#player2").onclick = function(){
      addPlayer.slot = 2;
      addPlayer.name = 'Achilles';
      addPlayer.color = 'red';
      socket.send(JSON.stringify(addPlayer));

      document.querySelector("#start").disabled = false;
      document.querySelector("#player1").style.display='none';
      document.querySelector("#player3").style.display='none';
      document.querySelector("#player4").style.display='none';
      document.querySelector("#player2").disabled = true;
      document.querySelector("#spectator").disabled = true;

  };

  document.querySelector("#player3").onclick = function(){
      addPlayer.slot = 3;
      addPlayer.name = 'Patrocle';
      addPlayer.color = 'blue';
      socket.send(JSON.stringify(addPlayer));
      document.querySelector("#start").disabled = false;
      document.querySelector("#player2").style.display='none';
      document.querySelector("#player1").style.display='none';
      document.querySelector("#player4").style.display='none';
      document.querySelector("#player3").disabled = true;
      document.querySelector("#spectator").disabled = true;

  };

  document.querySelector("#player4").onclick = function(){
      addPlayer.slot = 4;
      addPlayer.name = 'Ulysse';
      addPlayer.color = 'green';
      socket.send(JSON.stringify(addPlayer));
      document.querySelector("#start").disabled = false;
      document.querySelector("#player2").style.display='none';
      document.querySelector("#player3").style.display='none';
      document.querySelector("#player1").style.display='none';
      document.querySelector("#player4").disabled = true;
      document.querySelector("#spectator").disabled = true;

  };

/*  document.querySelector('#restart').onclick = function(){
      socket.send(JSON.stringify({kind : 'restart'}));
  };*/

  document.querySelector("#spectator").onclick = function(){
      addPlayer.slot = 0;
      addPlayer.name = 'Agamemnon';
      socket.send(JSON.stringify(addPlayer));
      document.querySelector("#start").disabled = false;
      document.querySelector("#spectator").disabled = true;
  };

  function checkSlot(slots){
      document.querySelector("#player2").disabled = true;
      document.querySelector("#player3").disabled = true;
      document.querySelector("#player1").disabled = true;
      document.querySelector("#player4").disabled = true;
     for(var i in slots){
          document.querySelector("#player" + slots[i]).disabled = false;
      }
  }


  function deadSnake(snake){
      document.getElementById("dead_message").innerHTML = "";
       for(var i in snake){
          if(snake[i].state == "dead"){
              document.getElementById("dead_message").innerHTML
              += snake[i].name + " is dead" + "<br>";
          }
      }
  }
};
