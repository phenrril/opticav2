<?php
require("../conexion.php");
session_start();

$porc_dto= $_POST['porc'];
$id_usuario = $_SESSION['idUser'];

    $query = mysqli_query($conexion, "INSERT INTO descuento(descuento, id_usuario)  VALUES ('$porc_dto', '$id_usuario')" );
    if($query){
        echo "<div class='alert alert-success'>DESCUENTO APLICADO CORRECTAMENTE</div>";
        echo '<script>document.getElementById("btn_descuento").setAttribute("type", "hidden");</script>';
        echo '<script>document.getElementById("btn_canceldto").setAttribute("type", "button");</script>';

    }
    else{
        echo "<div class='alert alert-success'>ERROR AL APLICAR DESCUENTO</div>";
    }

   // 

?>