<?php 
require "../conexion.php";

session_start();
$id = $_SESSION['idUser'];
$query2 = mysqli_query($conexion, "DELETE FROM descuento WHERE id_usuario = $id");

if($query2){
    echo "<div class='alert alert-success'>DESCUENTO ELIMINADO CORRECTAMENTE</div>";
    echo '<script>document.getElementById("btn_descuento").setAttribute("type", "button");</script>';
    echo '<script>document.getElementById("btn_canceldto").setAttribute("type", "hidden");</script>';
}
else{
    echo "<div class='alert alert-success'>ERROR AL ELIMINAR DESCUENTO</div>";
}


?>