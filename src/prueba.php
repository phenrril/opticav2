<?php 
require "../conexion.php";
include_once "includes/header.php";
$id_user = $_SESSION['idUser'];
$permiso = "nueva_venta";
$sql = mysqli_query($conexion, "SELECT p.*, d.* FROM permisos p INNER JOIN detalle_permisos d ON p.id = d.id_permiso WHERE d.id_usuario = $id_user AND p.nombre = '$permiso'");
$existe = mysqli_fetch_all($sql);
if (empty($existe) && $id_user != 1) {
    header("Location: permisos.php");
}




//cho "matias trolo";;
$sql3 = mysqli_query($conexion, "SELECT * from graduaciones where id_venta=57");

$fila = mysqli_fetch_array($sql3);
echo $fila['id'];

?>