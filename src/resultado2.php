<?php require "../conexion.php";


$fecha= $_POST['calendario'];


$sql=$conexion->query(" SELECT * FROM ventas  where date(fecha)='$fecha'");
        if ($datos=$sql->fetch_object()) {        
            while ($fila = mysqli_fetch_array($sql)){
            echo "<tr><td>".$fila['id']."</td><td>".$fila['id_cliente']."</td><td>".$fila['total']."</td><td>".$fila['id_usuario']."</td><td>".$fila['fecha']."</td></tr>";
                echo "<br />";
          }
            }












// $datos = array();
// $ventas = mysqli_query($conexion, "SELECT * FROM ventas where fecha='$fecha'");
//     while ($row = mysqli_fetch_assoc($ventas)) {
//         $data['id'] = $row['codproducto'];
//         $data['id_cliente'] = $row['codigo'] . ' - ' .$row['descripcion'];
//         $data['value'] = $row['descripcion'];
//         $data['precio'] = $row['precio'];
//         $data['existencia'] = $row['existencia'];
//         array_push($datos, $data);
//     }
//     echo json_encode($datos);




//echo $totalV;



?>